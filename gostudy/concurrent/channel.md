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
	//计算缓冲区总大小（elemsize*size）,判断是否超过最大可分配范围
	mem, overflow := math.MulUintptr(elem.size, uintptr(size))
	if overflow || mem > maxAlloc-hchanSize || size < 0 {
		panic(plainError("makechan: size out of range"))
	}

	//分配内存
	var c *hchan
	switch {
	case mem == 0:
		// 缓冲区大小为0，或者channel元素大小为0，仅分配channel必须的空间即可
		c = (*hchan)(mallocgc(hchanSize, nil, true))
		// Race detector uses this location for synchronization.
		c.buf = c.raceaddr()
	case elem.ptrdata == 0:
		// 元素不包含指针，分配一块连续内存
		// Allocate hchan and buf in one call.
		c = (*hchan)(mallocgc(hchanSize+mem, nil, true))
		c.buf = add(unsafe.Pointer(c), hchanSize)
	default:
		// 元素包含指针，为chan和缓冲区分别分配内存
		c = new(hchan)
		c.buf = mallocgc(mem, elem, true)
	}

	//设置属性
	c.elemsize = uint16(elem.size)
	c.elemtype = elem
	c.dataqsiz = uint(size)

	if debugChan {
		print("makechan: chan=", c, "; elemsize=", elem.size, "; dataqsiz=", size, "\n")
	}
	return c
}
```
#### 发送
通过sudog对g进行包装，以便携带数据项，存储相关状态
```go
type sudog struct {
	g *g
	next     *sudog
	prev     *sudog
	elem     unsafe.Pointer //elem 表示数据存储空间的指针
	c           *hchan
}
```

```go
func chansend1(c *hchan, elem unsafe.Pointer) {
	chansend(c, elem, true, getcallerpc())
}

func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
	//如果channel为nil
	if c == nil {
		//如果是非阻塞，直接返回发送不成功
		if !block {
			return false
		}
		//gopark会将当前goroutine阻塞挂起
		gopark(nil, nil, waitReasonChanSendNilChan, traceEvGoStop, 2)
		throw("unreachable")
	} 
	...
	//如果当前非阻塞且channel未关闭，如果无缓冲区且没有等待接收的Goroutine，或者有缓冲区且缓冲区已满，那么都直接返回发送不成功
	if !block && c.closed == 0 && full(c) {
		return false
	}
	//一些关于cpu的处理 参考https://github.com/golang/go/issues/8976
	var t0 int64
	if blockprofilerate > 0 {
		t0 = cputicks()
	}

	lock(&c.lock)
	//channel 已经关闭
	if c.closed != 0 {
		unlock(&c.lock)
		panic(plainError("send on closed channel"))
	}
	//查找当前channel中存在等待接收的goroutine
	if sg := c.recvq.dequeue(); sg != nil {
		//直接发送给接收的g，绕过缓冲区
		send(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true
	}
	//如果缓冲区没满
	if c.qcount < c.dataqsiz {
		// 通道缓冲区中有可用的空间。将要发送的元素排队。
		//qp是指向缓冲区sendx位置的指针
		qp := chanbuf(c, c.sendx)
		if raceenabled {
			raceacquire(qp)
			racerelease(qp)
		}
		//typedmemmove将要发送的值拷贝到缓冲区
		typedmemmove(c.elemtype, qp, ep)
		//因为是循环队列，sendx等于队列长度时置为0
		c.sendx++
		if c.sendx == c.dataqsiz {
			c.sendx = 0
		}
		c.qcount++
		unlock(&c.lock)
		return true
	}
	//当既没有等待接收的Goroutine，缓冲区也没有剩余空间，如果是非阻塞的发送，那么直接解锁，返回发送失败
	if !block {
		unlock(&c.lock)
		return false
	}

	//如果是阻塞发送，那么就将当前的Goroutine打包成一个sudog结构体，并加入到channel的发送队列sendq里
	gp := getg()
	mysg := acquireSudog()
	mysg.releasetime = 0
	if t0 != 0 {
		mysg.releasetime = -1
	}
	// No stack splits between assigning elem and enqueuing mysg
	// on gp.waiting where copystack can find it.
	mysg.elem = ep
	mysg.waitlink = nil
	mysg.g = gp
	mysg.isSelect = false
	mysg.c = c
	gp.waiting = mysg
	gp.param = nil
	c.sendq.enqueue(mysg)
	//调用goparkunlock将当前Goroutine设置为等待状态并解锁，进入休眠等待被唤醒
	gopark(chanparkcommit, unsafe.Pointer(&c.lock), waitReasonChanSend, traceEvGoBlockSend, 2)
	KeepAlive(ep)

	// 被唤醒之后执行清理工作并释放sudog结构体
	if mysg != gp.waiting {
		throw("G waiting list is corrupted")
	}
	gp.waiting = nil
	gp.activeStackChans = false
	if gp.param == nil {
		if c.closed == 0 {
			throw("chansend: spurious wakeup")
		}
		panic(plainError("send on closed channel"))
	}
	gp.param = nil
	if mysg.releasetime > 0 {
		blockevent(mysg.releasetime-t0, 2)
	}
	mysg.c = nil
	releaseSudog(mysg)
	return true
}

```
send 代码如下
```go
func send(c *hchan, sg *sudog, ep unsafe.Pointer, unlockf func(), skip int) {  
	  ...   
 if sg.elem != nil {        
   // 将发送的值直接拷贝到接收值（比如v = <-ch 中的v）的内存地址     
   sendDirect(c.elemtype, sg, ep)   
   sg.elem = nil  
  }    
  // 获取等待接收数据的Goroutine   
  gp := sg.g    
  unlockf()   
  gp.param = unsafe.Pointer(sg)  
  if sg.releasetime != 0 {      
	  sg.releasetime = cputicks()  
   }    
  // 唤醒之前等待接收数据的Goroutine   
	 
   goready(gp, skip+1)}
```
具体的发送过程
1. 首先会判断channel是否为nil
2. 判断channel的关闭情况
3. 根据缓冲区情况作出选择
   * 如果channel当前存在接受的goroutine，会直接将数据发送给该go程
   * 如果缓冲区没满，将正在发送的值拷贝到sendx位置，并增加sendx索引以及释放锁
   * 如果是阻塞发送，那么就将当前的Goroutine打包成一个sudog结构体，并加入到channel的发送队列sendq里。这个g会被设置为gwaiting表示正在阻塞，等待被唤醒

#### 接收
```go
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
	...
	//如果channel为nil，阻塞接受会将当前goroutine挂起
	if c == nil {
		if !block {
			return
		}
		gopark(nil, nil, waitReasonChanReceiveNilChan, traceEvGoStop, 2)
		throw("unreachable")
	}

	// Fast path: check for failed non-blocking operation without acquiring the lock.
	//如果非阻塞并且通道nil
	if !block && empty(c) {
		
		if atomic.Load(&c.closed) == 0 {
			// Because a channel cannot be reopened, the later observation of the channel
			// being not closed implies that it was also not closed at the moment of the
			// first observation. We behave as if we observed the channel at that moment
			// and report that the receive cannot proceed.
			return
		}
		// The channel is irreversibly closed. Re-check whether the channel has any pending data
		// to receive, which could have arrived between the empty and closed checks above.
		// Sequential consistency is also required here, when racing with such a send.
		if empty(c) {
			// The channel is irreversibly closed and empty.
			if raceenabled {
				raceacquire(c.raceaddr())
			}
			if ep != nil {
				typedmemclr(c.elemtype, ep)
			}
			return true, false
		}
	}

	var t0 int64
	if blockprofilerate > 0 {
		t0 = cputicks()
	}

	lock(&c.lock)
    // 如果channel已关闭，并且缓冲区无元素
	if c.closed != 0 && c.qcount == 0 {
		if raceenabled {
			raceacquire(c.raceaddr())
		}
		unlock(&c.lock)
		 // 有等待接收的变量（即 v = <-ch中的v）
		if ep != nil {
			 //根据channel元素的类型清理ep对应地址的内存，即ep接收了channel元素类型的零值
			typedmemclr(c.elemtype, ep)
		}
		return true, false
	}
	//对于发送队列里如果也有等待的goroutine，两种情况要么channel满了要么channel没有缓冲区
	if sg := c.sendq.dequeue(); sg != nil {
		// Found a waiting sender. If buffer is size 0, receive value
		// directly from sender. Otherwise, receive from head of queue
		// and add sender's value to the tail of the queue (both map to
		// the same buffer slot because the queue is full).
		//如果无缓冲区，那么直接从sender接收数据；否则，从buf队列的头部接收数据，并把sender的数据加到buf队列的尾部
		recv(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true, true
	}
	//如果队列里面有数据
	if c.qcount > 0 {
		// Receive directly from queue
		//从recvx指向的位置获取元素
		qp := chanbuf(c, c.recvx)
		if raceenabled {
			raceacquire(qp)
			racerelease(qp)
		}
		if ep != nil {
			//将从buf中取出的元素拷贝到当前协程
			typedmemmove(c.elemtype, ep, qp)
		}
		//同时将取出的数据所在的内存清空
		typedmemclr(c.elemtype, qp)
		c.recvx++
		if c.recvx == c.dataqsiz {
			c.recvx = 0
		}
		c.qcount--
		unlock(&c.lock)
		return true, true
	}

	if !block {
		unlock(&c.lock)
		return false, false
	}

	// no sender available: block on this channel.
	//阻塞模式，获取当前Goroutine，打包一个sudog
	gp := getg()
	mysg := acquireSudog()
	mysg.releasetime = 0
	if t0 != 0 {
		mysg.releasetime = -1
	}
	// No stack splits between assigning elem and enqueuing mysg
	// on gp.waiting where copystack can find it.
	mysg.elem = ep
	mysg.waitlink = nil
	gp.waiting = mysg
	mysg.g = gp
	mysg.isSelect = false
	mysg.c = c
	gp.param = nil
	//// 加入到channel的等待接收队列recvq中
	c.recvq.enqueue(mysg)
	gopark(chanparkcommit, unsafe.Pointer(&c.lock), waitReasonChanReceive, traceEvGoBlockRecv, 2)

	// someone woke us up
	//被唤醒之后执行清理工作并释放sudog结构体
	if mysg != gp.waiting {
		throw("G waiting list is corrupted")
	}
	gp.waiting = nil
	gp.activeStackChans = false
	if mysg.releasetime > 0 {
		blockevent(mysg.releasetime-t0, 2)
	}
	closed := gp.param == nil
	gp.param = nil
	mysg.c = nil
	releaseSudog(mysg)
	return true, !closed
}
```
1. 如果sendq存在挂起的goroutine，调用recv()
   1.  如果无缓冲区，那么直接从sender接收数据；
   2.  如果缓冲区已满，从buf队列的头部接收数据，并把sender的数据加到buf队列的尾部；
   3.  最后调用goready函数将等待发送数据的Goroutine的状态从_Gwaiting置为_Grunnable，等待下一次调度。
2. 如果缓冲区还存在数据，将从buf中取出的元素拷贝到当前协程的接收数据目标内存地址中。
3. 如果是阻塞模式，且当前没有数据可以接收，那么就需要将当前Goroutine打包成一个sudog加入到channel的等待接收队列recvq中，将当前Goroutine的状态置为_Gwaiting，等待唤醒
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

