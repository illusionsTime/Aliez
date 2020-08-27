## channel

### 介绍 
Go鼓励“应该以通信作为手段来共享内存”，channel则是最直接的体现。channel提供了一种机制，它既可以同步两个并发执行的函数，又可以让这两个函数通过互相传递特定类型的值来通信。

数据传递方向的不同意味着他们类型的不同。
*同步与异步的区别在于是否带有缓冲槽* 调用内建函数cap判断

* FIFO
目前的 Channel 收发操作均遵循了先入先出（FIFO）的设计，具体规则如下：

1. 先从 Channel 读取数据的 Goroutine 会先接收到数据；
2. 先向 Channel 发送数据的 Goroutine 会得到先发送数据的权利；

#### 数据结构
runtime.go/chan
```go
type hchan struct {
	qcount   uint           // 队列中的总数据
	dataqsiz uint           // 缓冲槽大小
	buf      unsafe.Pointer // 指向dataqsiz的指针
	elemsize uint16         // 数据项大小
	closed   uint32         
	elemtype *_type // 数据项类型
	sendx    uint   // Channel 的发送操作处理到的位置；
	recvx    uint   // Channel 的接收操作处理到的位置；
	recvq    waitq  // list of recv waiters
	sendq    waitq  // list of send waiters
	lock mutex
}
```
sendq 和 recvq 存储了当前 Channel 由于缓冲区空间不足而阻塞的 Goroutine 列表，这些等待队列使用双向链表 runtime.waitq 表示，
```go
type waitq struct {
	first *sudog
	last  *sudog
}
```
runtime.sudog 表示一个在等待列表中的 Goroutine，该结构体中存储了阻塞的相关信息以及两个分别指向前后 runtime.sudog 的指针。
#### 创建
golang中所有channel的创建都会遵循make关键字
chan.go
```go
func makechan(t *chantype, size int) *hchan {
	elem := t.elem

	// compiler checks this but be safe.
	if elem.size >= 1<<16 {
		throw("makechan: invalid channel element type")
	}
	if hchanSize%maxAlign != 0 || elem.align > maxAlign {
		throw("makechan: bad alignment")
	}

	mem, overflow := math.MulUintptr(elem.size, uintptr(size))
	if overflow || mem > maxAlloc-hchanSize || size < 0 {
		panic(plainError("makechan: size out of range"))
	}

	// Hchan does not contain pointers interesting for GC when elements stored in buf do not contain pointers.
	// buf points into the same allocation, elemtype is persistent.
	// SudoG's are referenced from their owning thread so they can't be collected.
	// TODO(dvyukov,rlh): Rethink when collector can move allocated objects.
	var c *hchan
	switch {
	case mem == 0:
		// Queue or element size is zero.
		c = (*hchan)(mallocgc(hchanSize, nil, true))
		// Race detector uses this location for synchronization.
		c.buf = c.raceaddr()
	case elem.ptrdata == 0:
		// Elements do not contain pointers.
		// Allocate hchan and buf in one call.
		c = (*hchan)(mallocgc(hchanSize+mem, nil, true))
		c.buf = add(unsafe.Pointer(c), hchanSize)
	default:
		// Elements contain pointers.
		c = new(hchan)
		c.buf = mallocgc(mem, elem, true)
	}

	c.elemsize = uint16(elem.size)
	c.elemtype = elem
	c.dataqsiz = uint(size)

	if debugChan {
		print("makechan: chan=", c, "; elemsize=", elem.size, "; dataqsiz=", size, "\n")
	}
	return c
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

