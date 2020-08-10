### 并发和并行的概念

并行：在任意粒度的时间内都具备同时执行的能力
并发：程序在单位时间内是同时运行的

### 多进程编程

Linux中支持IPC的方法有很多种。从处理机制看 可分为基于通信的IPC方法、基于信号的IPC方法、基于同步的IPC方法。
其中，基于通信的IPC方法又分为以数据传送为手段的IPC方法和一共享内存为手段的IPC方法，前者包含了管道和消息队列。
管道可以用来传送字节流、消息队列可以用来传送结构化的消息对象。以共享内存为手段的IPC方法主要以共享内存区为代表，他是最快的一种IPC方法。
基于信号的IPC方法就是我们常说的os的信号机制，他是惟一异步的IPC方法。

*Go中支持IPC的方法有管道、信号、和socket*

### 进程和线程以及协程

当操作系统运行一个应用程序的时候，os会为这个程序启动一个进程。可以说这个进程是一个包含了应用程序在运行中需要用到和维护的各种资源的容器
* 进程是具有一定独立功能的程序关于某个数据集合上的一次运行活动,进程是系统进行资源分配和调度的一个独立单位。每个进程都有自己的独立内存空间，不同进程通过进程间通信来通信。由于进程比较重量，占据独立的内存，所以上下文进程间的切换开销（栈、寄存器、虚拟内存、文件句柄等）比较大，但相对比较稳定安全

* 线程 线程是进程的一个实体,是CPU调度和分派的基本单位,它是比进程更小的能独立运行的基本单位.线程自己基本上不拥有系统资源,只拥有一点在运行中必不可少的资源(如程序计数器,一组寄存器和栈),但是它可与同属一个进程的其他的线程共享进程所拥有的全部资源。线程间通信主要通过共享内存，上下文切换很快，资源开销较少，但相比进程不够稳定容易丢失数据。

* 协程是一种用户态的轻量级线程，协程的调度完全由用户控制。协程拥有自己的寄存器上下文和栈。协程调度切换时，将寄存器上下文和栈保存到其他地方，在切回来的时候，恢复先前保存的寄存器上下文和栈，直接操作栈则基本没有内核切换的开销，可以不加锁的访问全局变量，所以上下文的切换非常快。

<font color=red>为什么说协程的执行效率高？</font>
线程的切换时os来进行的，涉及模式转换需要在内核态和用户态间来回切换，一次切换就需要它的寄存器和内核栈进行刷新
协程的切换是用户自己来决定的，切换不设计模式转换，仅用户态，只需要对自己拥有的寄存器的只进行修改


* 进程的状态

1. 可运行状态 R
2. 可中断的睡眠状态 S
3. 不可中断的睡眠状态 D
4. 暂停状态或者跟踪状态 T
5. 僵尸状态 Z  处于此状态的进程即将结束运行 该进程占用的绝大多数资源也已经回收 ，不过还有一些信息未删除，比如退出码以及一些统计信息。之所以保留
   这些信息，主要是考虑到该进程的父进程可能需要他们。
6. 退出状态 X

#### 系统调用
用户进程生存在用户空间，不能与计算机硬件进行交互。内核可以与硬件交互，但是却生存在内核空间。
为了使用户进程能够使用操作系统更底层的功能，内核会暴露出一些接口供他们使用，这些接口时用户能够使用内核功能的唯一手段。
用户进程使用这些接口的行为称为系统调用。

* 内核态和用户态
  为了保证操作系统的稳定和安全，内核依据由CPU提供的，可以让进程驻留的特权级别建立了两个特权状态——内核态和用户态。
  大部分时间，CPU都处于用户态，这是CPU只能对用户空间进行访问。也就是说CPU在用户态下运行的用户进程是不能与内核接触的。当用户发出一个系统调用时，内核会把CPU从用户态切换到内核态，而后让CPU执行对应的内核函数。CPU在内核态下是有权访问内核空间的。当内核函数执行完毕后，内核会把CPU从内核态切换回用户态，并把执行结果返回给用户进程。