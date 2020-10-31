# Pprof
Pprof可以进行golang的程序的性能监控及分析

可以做什么：
1. CPU分析 Pprof可以按一定频率监听程序消耗CPU的情况，可确定应用程序在主动消耗CPU周期是所处的位置
2. 内存分析 在应用程序进行堆分配时记录堆栈跟踪 以便检查内存泄露
3. 阻塞分析 记录goroutine阻塞等待同步的位置
4. 互斥锁分析 报告互斥锁的竞争情况

## 分析
### Web界面
* cpu（CPU Profiling）: $HOST/debug/pprof/profile，默认进行 30s 的 CPU Profiling，得到一个分析用的 profile 文件
* block（Block Profiling）：$HOST/debug/pprof/block，查看导致阻塞同步的堆栈跟踪
* goroutine：$HOST/debug/pprof/goroutine，查看当前所有运行的 goroutines 堆栈跟踪
* heap（Memory Profiling）: $HOST/debug/pprof/heap，查看活动对象的内存分配情况
* mutex（Mutex Profiling）：$HOST/debug/pprof/mutex，查看导致互斥锁的竞争持有者的堆栈跟踪
* threadcreate：$HOST/debug/pprof/threadcreate，查看创建新OS线程的堆栈跟踪

### 终端
flat：给定函数上运行耗时
flat%：同上的 CPU 运行耗时总比例
sum%：给定函数累积使用 CPU 总比例
cum：当前函数加上它之上的调用运行总耗时
cum%：同上的 CPU 运行耗时总比例

