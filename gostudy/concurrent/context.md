context :退出通知 ，元数据传递

设计目地：跟踪goroutine调用，在其内部维护一个调用，并在这些调用树中传递通知和元数据。

核心功能：多个goroutine之间的退出通知机制。

#### context接口


empty Ctx 结构
：作为context对象树的根，