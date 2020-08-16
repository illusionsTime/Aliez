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

#### 规则 
在开始执行select语句的时候，所有跟在case右边的发送语句或者接受语句中的通道表达式和元素表达式都会先求值(求值的谁徐是从左到右，自上而下的)

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

