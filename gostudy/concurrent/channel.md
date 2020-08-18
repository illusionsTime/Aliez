## channel

#### 介绍 
Go鼓励“应该以通信作为手段来共享内存”，channel则是最直接的体现。channel提供了一种机制，它既可以同步两个并发执行的函数，又可以让这两个函数通过互相传递特定类型的值来通信。

数据传递方向的不同意味着他们类型的不同。
*同步与异步的区别在于是否带有缓冲槽* 调用内建函数cap判断

runtime.go/chan
```go
type hchan struct {
	qcount   uint           // 队列中的总数据
	dataqsiz uint           // 缓冲槽大小
	buf      unsafe.Pointer // 指向dataqsiz的指针
	elemsize uint16         // 数据项大小
	closed   uint32         
	elemtype *_type // 数据项类型
	sendx    uint   // send index
	recvx    uint   // receive index
	recvq    waitq  // list of recv waiters
	sendq    waitq  // list of send waiters
	lock mutex
}
```

### select 
select语句是一种仅能用于通道发送和接受操作得专用语句。一条select语句执行时会选择其中的某一条分支并执行。
```go
var intChan=make(chan int,10)
var strChan=make(chan string,10)
select {
   case e1:=<-intChan:
        fmt.Printf("The 1th case was selected. e1=%v.\n",e1)
   case e2:=<-strChan:
        fmt.Printf("The 2nd case was selected. e2=%v.\n",e2)
   default:
   fmt.Println("default!")
}
```


#### 分支选择规则 
在开始执行select语句的时候，<font color=red>所有跟在case右边的发送语句或者接受语句中的通道表达式和元素表达式都会先求值(求值的顺序是从左到右，自上而下的)，无论所有case是否有可能被选择都会这样</font>

在执行select语句时，运行时系统会自上而下地判断每个case中的发送或接收操作是否可以立即执行。这里的立即执行指的是当前goroutine不会因此操作而被阻塞。只要发现有一个case上的判断是肯定的，该case就会被选中。

当有一个case被选中时，运行时系统就会执行该case及其包含的语句，而其他case被忽略。如果有多个case满足条件，那么运行时系统会通过一个*伪随机*的算法选中一个case。例如channel.go中提到的，多次运行的结果都会不尽相同。
`随机的引入避免了饥饿问题的发生`

channel.go
```go
func pse() {
	chanCap := 5
	intChan := make(chan int, chanCap)
	for i := 0; i < chanCap; i++ {
		select {
		case intChan <- 1:
		case intChan <- 2:
		case intChan <- 3:
		}
	}
	for i := 0; i < chanCap; i++ {
		fmt.Printf("%d\n", <-intChan)
	}
}

```

* 如果select语句中所有的case都不满足选择条件并且没有default case,那么当前goroutine就会一直阻塞于此，直到至少有一个case中的发送或者接受操作可以立即执行为止。

### future模式
#### 工作原理
1. 使用chan作为函数参数
2. 启动goroutine调用函数
3. 通过chan传入参数
4. 做其他可以并行处理的事情
5. 通过chan异步获取结果

```
+----+
|Main|                                                                   task                  Result
+----+                                                                   chan                   chan
   |     go      +-----------+                                            ||                     ||
   |------------>|future task|                                            ||                     ||
   |             +-----------+                                            ||                     ||
   |                  |                      Read task                    ||                     ||
   |                  |<--------------------------------------------------||                     ||
   |                  |                     Write Ressult                 ||                     ||
   |                  |---------------------------------------------------||-------------------->||
   |                                                                      ||                     ||
   |                 Write task                                           ||                     ||
   |--------------------------------------------------------------------->||                     ||
   |                                                                                             ||
   |                                                                                             ||
{Do Something}                                                                                   ||
   |                                                                                             ||
   |                           Read   Result                                                     ||
   |<--------------------------------------------------------------------------------------------||
   |
   |
```

