<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [滴滴（自动驾驶Vovager）](#滴滴自动驾驶vovager)
	- [一面](#一面)
	- [二面](#二面)
	- [三面](#三面)
- [富途](#富途)
- [最右](#最右)
	- [一面](#一面-1)
- [百度 （AI平台）](#百度-ai平台)
	- [一面](#一面-2)
	- [二面](#二面-1)
	- [三面 （现场）](#三面-现场)
	- [四面  （现场）](#四面-现场)
- [b站](#b站)
	- [一面](#一面-3)
	- [二面](#二面-2)
- [Aibee](#aibee)
	- [一面](#一面-4)
	- [二面](#二面-3)
	- [三面](#三面-1)
- [滴滴 （网约车地图团队）](#滴滴-网约车地图团队)
	- [一面](#一面-5)
- [小米 Miot](#小米-miot)
	- [一面](#一面-6)
	- [二面](#二面-4)
- [字节跳动 懂车帝](#字节跳动-懂车帝)
	- [一面](#一面-7)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 滴滴（自动驾驶Vovager）
### 一面
主要是项目以及部分基础知识，算法的话是股票问题，还是很简单的动态规划问题  
剑指offer.63

假设把某股票的价格按照时间先后顺序存储在数组中，请问买卖该股票一次可能获得的最大利润是多少？

示例 1:
```
输入: [7,1,5,3,6,4]
输出: 5
```
解释: 在第 2 天（股票价格 = 1）的时候买入，在第 5 天（股票价格 = 6）的时候卖出，最大利润 = 6-1 = 5 。  
     注意利润不能是 7-1 = 6, 因为卖出价格需要大于买入价格。  

思路：在第i天时，有两种选择,不出售股票和出售股票，那么不出售股票的利润就是前一天的利润dp[i-1],  
出售的话期利润就是num[i]-num[j],在第j天购买股票，如何判断j是哪一天我们通过一个常量来记录前面遍历过的最小值。

可以写出状态转移方程dp[i]=max(dp[i-1],num[i]-num[j])

```go
func maxProfit(prices []int) int {
	size := len(prices)
    if size==0||size==1{
        return 0
    }
	profit := 0
	cost := prices[0]
	for i := 1; i < size; i++ {
		cost = min(prices[i], cost)
		profit = max(profit, prices[i]-cost)
	}
	return profit
}
```
* sql题 查找当前班级所有科目平均分大于90的学生姓名  
``` sql 
select name from table group by name having avg(score)>90;
```

### 二面 
二面主要问了项目以及一道算法题和场景设计题，我个人觉得不具有总结性

### 三面
三面的话感觉状态太差，一道送分题变成送命题，后面面试官也就没问什么就结束了

剑指offer 29 顺时针打印矩阵
这道题应该采用暴力的按层模拟打印即可，需要注意的是打印的边界。

![](../views/顺时针打印矩阵.png)

```go
func spiralOrder(matrix [][]int) []int {
	if matrix == nil || len(matrix) == 0 || len(matrix) == 0 {
		return []int{}
	}
	top := 0
	hsize := len(matrix)
	lsize := len(matrix[0])
	left := 0

	bottom := hsize - 1
	right := lsize - 1
	index := 0
	x, y := 0, 0
	sum := make([]int, hsize*lsize)
	for bottom >= top && right >= left {
		for x = left; x <= right; x++ {
			sum[index] = matrix[top][x]
			index++
		}
		for y = top + 1; y <= bottom; y++ {
			sum[index] = matrix[y][right]
			index++
		}
		if bottom > top && right > left {
			for x = right - 1; x > left; x-- {
				sum[index] = matrix[bottom][x]
				index++
			}
			for y = bottom; y > top; y-- {
				sum[index] = matrix[y][left]
				index++
			}
		}
		left++
		right--
		top++
		bottom--
	}
	return sum
}
```

## 富途
主要是三道算法题 第一二道要写出来，第三道讲思路
先讲第三道吧  
设计一个算法，找出数组中最小的k个数。以任意顺序返回这k个数均可，时间复杂度O(nlogk) 

个人首先想到的是堆排 构建一个小顶堆，每次将堆顶元素滞后，重新下滤，这样得到的是排序后的k个数

第二种方法是快排的变种，找出枢纽元的index与k的关系，相等则返回，index>k证明在枢纽元前边，反之则在后边，时间复杂度应该是等比数列的求和O(n),但是存在最坏情况O(n^2),对于k较大时，这应该是最高效的解法。

```go
func getLeastNumbers(arr []int, k int) []int {
	if k >= len(arr) {
		return arr
	}
	return quickselect(arr, 0, len(arr)-1, k)
}

func quickselect(arr []int, left, right int, k int) []int {
	if left < right {
		index := partition(arr, left, right)
		if index == k {
			return arr[:k]
		} else if index > k {
			return quickselect(arr, left, index-1, k)
		} else if index < k {
			return quickselect(arr, index+1, right, k)
		}
	}
	return arr[:k]
}

func partition(num []int, l, r int) int {
	poivt := num[l]
	index := l + 1
	for i := index; i <= r; i++ {
		if num[i] < poivt {
			num[i], num[index] = num[index], num[i]
			index++
		}
	}
	num[l], num[index-1] = num[index-1], num[l]
	return index - 1
}

```

第三种方法是如果数组范围极为有限，通过一个桶即可实现一次遍历完成。

第二道算法题 leetcode209
给定一个含有 n 个正整数的数组和一个正整数 s ，找出该数组中满足其和 ≥ s 的长度最小的 连续 子数组，并返回其长度。如果不存在符合条件的子数组，返回 0。

示例：
```
输入：s = 7, nums = [2,3,1,2,4,3]
输出：2
```
解释：子数组 [4,3] 是该条件下的长度最小的子数组。

思路 暴力法需要O(n^2)，想更进一步可以用二分查找，最好的方法是双指针法时间复杂度为O(n)

```go
func minSubArrayLen(s int, nums []int) int {
	size := len(nums)
	if size == 0 {
		return 0
	}
	sum := 0
	ans := size + 1
	start, end := 0, 0
	for end < size {
		sum += nums[end]
		for sum >= s {
			ans = min(ans, end-start+1)
			sum -= nums[start]
			start++
		}
		end++
	}
	if ans == size+1 {
		return 0
	}
	return ans
}
```

## 最右 
### 一面
* channel在生产端还是消费端关闭？
  *不要在消费端关闭channel，不要在有多个并行的生产者时对channel执行关闭操作。*
  如果要在消费端关闭呢？  
  答：可以使用recover机制，避免程序因为panic而崩溃
```go
  func SafeClose(ch chan T) (justClosed bool) {
	defer func() {
		if recover() != nil {
			justClosed = false
		}
	}()
	
	// assume ch != nil here.
	close(ch) // panic if ch is closed
	return true // <=> justClosed = true; return
  }
 ```
如何优雅地关闭channel?
假如是多个消费者，单个生产者，直接在生产者端关闭

如果是多个生产者，单个消费者，可以通过一个通知关闭channel告知生产者，也就是说生产者同时也是告知信号的接受者，退出信号channel仍然是由它的生产端关闭的，所以这仍然没有违背channel关闭原则。

如果是多个生产者多个消费者，可以引入一个额外的协调者来关闭附加的退出信号channel。



* map多个读会发生什么？怎么实现线程安全？
  golang的map不是并发安全的，假入并发量不是很大的情况下map不会出现问题，但是当并发量很大的情况下会出现error并发写map异常，读的话没有问题，想要实现线程安全可以通过sync.RWMutex的方式实现

* new和make的区别
new和make都在堆上分配内存，但是它们的行为不同，适用于不同的类型。

new(T) 为每个新的类型T分配一片内存，初始化为 0 并且返回类型为*T的内存地址：这种方法 返回一个指向类型为 T，值为 0 的地址的指针，它适用于值类型如数组和结构体；它相当于 &T{}。

make的目的不同于new，它只用于slice,map,channel的创建，并返回类型为T（非指针）的已初始化（非零值）的值；出现这种差异的原因在于，这三种类型本质上为引用类型，它们在使用前必须初始化；


* Redis zset的实现
  REDIS_ZSET类型对应两种编码REDIS_ENCODING_SKIPLIST以及REDIS_ENCODING_ZIPLIST，REDIS_ENCODING_ZIPLIST对应的底层数据结构是压缩列表，REDIS_ENCODING_SKIPLIST使用的数据结构是跳跃表和字典。

  当有序集合对象可以同时满足两个条件时，对象使用ziplist编码
   1. 有序集合保存的元素数量小于128个
   2. 有序集合保存的所有元素成员的长度都小于64字节
  不能满足以上两个条件的有序集合对象将使用skiplist编码

ziplist编码的压缩列表对象使用压缩列表作为底层实现，每个集合元素使用两个紧挨在一起的压缩列表节点来保存，第一个节点保存元素的成员，第二个节点保存元素的分值。压缩列表内的集合元素按分值从小到大进行排序，分值较小的元素被放置在表头方向。

如果是skiplist编码的有序集合对象使用zset作为底层实现，一个zset结构是包含一个字段和一个跳跃表
```c
typedef struct zset{
	zskiplist *zsl;
	dict *dict;
}zset
```
zset结构中的zsl跳跃表按分值从小到大保存了所有的集合元素，每个跳跃表节点保存一个集合元素；字典的键保存了元素的成员，而字典的值保存了元素的分值。通过这个字典，程序可以O(1)复杂度超找给定元素的分值。*值得一提的是，虽然zset结构通式使用跳跃表和字典来保存有序集合元素，但是这两种数据结构都会通过指针来共享相同元素的成员和分值，所以同时使用跳跃表和字典来保存集合元素不会产生任何重复成员或者分值，也不会浪费额外的内存*

为什么需要同时使用跳跃表和字典来实现？
因为单独使用跳跃表或者字典在性能上比起使用跳跃表和字典都会有所降低，如果只是用字典，虽然能以O（1）复杂度查找成员的分值这一特性会被保留，但是因为字典以无序的方式来保存集合元素，所以在执行范围型操作时，需要现对字典中的元素进行排序，完成这种排序需要至少O(NlogN)时间复杂度，以及额外的O(N)内存空间。  
而如果单独使用跳跃表，那么跳跃表执行范围型操作的优点会被保留，但是根据成员查找分值这一操作的复杂度将上升为O(logN)

## 百度 （AI平台）
### 一面
* redis和etcd的区别？
1. 从数据结构方面来讲 Redis支持多种数据类型（string，set，list，hash，zset）
2. 从读写性能上来讲，Redis读写性能优异，并且提供了RDB、AOF持久化，而etcd v3的底层采用boltdb做存储，value直接持久化
3. 从使用场景上来看，etcd更适用于服务发现，配置管理，而Redis更适用于非强一致性的需求，比如说是队列，缓存，分布式Session
4. 两者都是KV存储，但是etcd通过Raft算法保证了各个节点间的数据和事务的一致性，更强调各个节点间的通信；Redis则时更像是内存式的缓存，因此来说读写能力很强。
5. Redis是c开发的，etcd是go开发的，他是源于k8s的兴起作为一个服务发现。
6. etcd v3只能通过gRPC访问，而redis可以通过http访问，因此etcd的客户端开发工作量高很多。

* etcd的Raft算法介绍一下
  常见的可靠一致性算法有Paxos 和 Raft，Raft算法的具体介绍在etcd里做了介绍，就不再多讲了，因为可能会被问到如果能答上来会加分很多吧

* ping指令的实现
这个在tcp/ip详解里是有的，可惜偷懒没认真看。


* linux的用户态和内核态
  为了保证操作系统的稳定和安全，内核依据由CPU提供的，可以让进程驻留的特权级别建立了两个特权状态——内核态和用户态。
  大部分时间，CPU都处于用户态，这是CPU只能对用户空间进行访问。也就是说CPU在用户态下运行的用户进程是不能与内核接触的。当用户发出一个系统调用时，内核会把CPU从用户态切换到内核态，而后让CPU执行对应的内核函数。CPU在内核态下是有权访问内核空间的。当内核函数执行完毕后，内核会把CPU从内核态切换回用户态，并把执行结果返回给用户进程。
* 什么时候从用户态切换为内核态
  系统调用，异常如缺页异常，外部设备中断
* 为什么说线程的切换比协程慢
  主要还是一个上下文切换问题，协程只有用户态，在切换时不涉及模式转换，只是把当前协程的 CPU 寄存器状态保存起来，然后将需要切换进来的协程的 CPU 寄存器状态加载的 CPU 寄存器上就 ok 了。而且完全在用户态进行
  线程切换的话首先除了和协程相同基本的 CPU 上下文，还有线程私有的栈和寄存器等，其次线程切换在切换时，先保存CPU寄存器中用户态的指令位置，再重新更新为内核指令的位置。当系统调用结束时，CPU寄存器恢复到原来保存的用户态。一次系统调用，发生了两次CPU上下文切换。
* redis/etcd对分布式锁的实现
  大致概念一样，都是通过key的创建作为一个互斥量，如果这个key不存在，那么线程创建key，成功则获取到锁，key为存在状态；
  如果key已经存在，那么线程不能创建key，获取锁失败。

  
### 二面
  百度的面试体验真的是面试过最好的，面试官会根据你的项目设计题目，和你一起探讨解决方法。
  最后出了一道算法题，有10个元素存到长度为12的数组中，有两个元素重复，找出这两个元素
  开始想到通过一次遍历异或的过程，但是只限于一个元素重复，后面面试官提醒了一下可以找到ab元素的关系  
  比如a+b通过把元素累加在进行减法运算可以得到a+b
  通过累乘再进行除法运算可以得到a*b的值  
  考虑到假如元素溢出，那么可已通过平方累加得到a^2+b^2的值  
  根据关系式可以求得a，b的值

### 三面 （现场）
问了go的协程，调度，包括cpu上下文切换的问题
项目挖了一些内容介绍等等
三道code，堆排，string化整型，反转链表
### 四面  （现场）
问了一个智力题，一个求单链表中间节点（快慢指针），一道dp题目

然后问了一些项目的问题，遇到什么该怎么解决，如何查找代码的性能瓶颈等等

然后没有hr，把我送下楼了，可能凉了吧
## b站
### 一面
* go和java的区别？
  1. java允许多态（重载和覆盖），go没有，同样的java支持函数重载，golang中函数必须具有唯一性
  2. 对于运行速度来说，go比java的执行效率要快
  3. go对于并发的支持比java要好，goroutine和channel
  4. Go语言的继承通过匿名组合完成：基类以Struct的方式定义，子类只需要把基类作为成员放在子类的定义中，支持多继承。Java 中通过 extends 关键字可以申明一个类是从另外一个类继承而来的
  5. golang的面向对象主要通过struct和interface来实现，go中没有类的概念，go是面向接口编程。
  6. go的包管理和java的包管理
   
* mysql事务隔离级别？
  脏读： 该隔离级别的事务会读到其它未提交事务的数据，此现象也称之为 脏读 。
  不可重复读： 一个事务可以读取另一个已提交的事务，多次读取会造成不一样的结果，此现象称为不可重复读问题
  幻读： 是针对数据插入（INSERT）操作来说的。假设事务A对某些行的内容作了更改，但是还未提交，此时事务B插入了与事务A更改前的记录相同的记录行，并且在事务A提交之前先提交了，而这时，在事务A中查询到了记录行未更改
  1. READ UNCOMMITTED 读未提交 都可能发生
  2. READ COMMITTED   解决了脏读问题
  3. REPEATABLE READ 可重复读 解决了脏读不可重复读问题，可能会发生幻读 InnoDB支持的隔离级别，但是使用NEXT-key算法避免了幻读的产生 
  4. SERIALIZABLE 序列化 在该隔离级别下事务都是串行顺序执行的，MySQL 数据库的 InnoDB 引擎会给读操作隐式加一把读共享锁，从而避免了脏读、不可重读复读和幻读问题。
   
* tcp四次挥手？为什么要四次挥手？
* 一道sql题，tb1: id,name  tb2:id,num 找出成绩最好的前100个人的名字
  select name from tb1,tb2 where tb2.num>=( select min (num) from (SELECT top 100 num from tb2
order by num))and tb1.id==tb2.id
* 一道golang的闭包问题
  参考之前写过的闭包问题
* golang并发模型
  Golang 就是借用CSP模型的一些概念为之实现并发进行理论支持，其实从实际上出发，go语言并没有，完全实现了CSP模型的所有理论，仅仅是借用了 process和channel这两个概念。process是在go语言上的表现就是 goroutine 是实际并发执行的实体，每个实体之间是通过channel通讯来实现数据共享。

* channel有缓冲槽和无缓冲槽的区别？
* 负载均衡的实现？

负载均衡：负载均衡是指，将请求分发到 多台 应用服务器，以此来分散 压力的一种架构方式，他是以集群的方式存在，并且当 某个节点挂掉的时候，可以自动 不再将请求分配到此节点。

1. 重定向 是通过将请求全部发送到前置机，由前置机通过算法 得出要分配给那台 应用服务器，然后响应给客户端，由客户端重定向到应用服务器的一种方式。 *效率低*
2. 反向代理 是通过在前置机，使用反向代理的方式，将请求分发到应用服务器，客户端无需再请求一次，实现方式通常有两种，一种是用交换机实现，还有一种是用nginx这一类的软件实现 *效率较高但是压力大*
3. 数据链路层处理 通过给应用服务器设置虚拟IP，然后通过修改mac地址的方式，将请求分发出去，而应用服务器 收到请求后，可以直接响应给客户端，而不需要经过前置机。
### 二面
## Aibee
### 一面
* 数组中找出和为k的两个元素返回下标
  双指针法：前提是数组是递增的  
```go
func twoSum(nums []int, target int) []int {
	if len(nums) < 2 {
		return nil
	}
	i, j := 0, len(nums)-1
	for i != j {
		if nums[i]+nums[j] == target {
			return []int{nums[i], nums[j]}
		} else if nums[i]+nums[j] < target {
			i++
		} else if nums[i]+nums[j] > target {
			j--
		}
	}
	return nil
}
```
另一道有意思的题目可以参考leetcode 560题


* 无头结点的链表元素删除

* https的加密过程
  首先client请求server端，server返回一个证书公钥，client会验证证书的合法性，随后生成一个随机值并通过公钥加密随机值，将加密后的密钥发给server端，server会通过私钥解密密钥，然后用密钥加密要发送的内容给client，client会通过密钥解析。

* 排行榜的设计
  这个主要是基于Redis的zset的实现
  ZAdd/ZRem是O(log(N))，ZRangeByScore/ZRemRangeByScore是O(log(N)+M)，N是Set大小，M是结果/操作元素的个数。
  ZSET的实现用到了两个数据结构：hash table 和 skip list(跳跃表)，其中hash table是具体使用redis中的dict来实现的，主要是为了保证查询效率为O(1) ，而skip list(跳跃表)主要是保证元素有序并能够保证INSERT和REMOVE操作是O(logn)的复杂度。

* InnoDb的索引模型
  哈希索引 b+索引 全文索引
  
### 二面
二面主要是一道算法题和股票差不多，包括项目和slice的一些问题。

### 三面
和面试官互相吹牛逼了

## 滴滴 （网约车地图团队）
### 一面
问了很多项目包括golang的一些基础实现以及源码

* redis的string类型的底层实现 和c的区别
1. 获取字符串长度只需要O（1）的时间复杂度，程序仅需要访问SDS的len属性即可
2. c语言会产生缓冲区溢出，在Redis中如果需要对SDS进行修改时，API回显检查SDS的空间是否满足修改所需要的要求，如果不满足的话，API会自动将SDS的空间拓展到执行修改所需要的大小，然后才执行实际的修改操作
3. 减少修改字符串时带来的内存重分配次数，SDS通过free属性记录未使用的字节，通过未使用空间，SDS实现了空间预分配和惰性空间释放两种优化策略。是的修改字符串长度N次最多需要执行N次的内存重分配。
4. c字符串只能保存文本数据，SDS可以保存文本或者二进制数据
  
* redis常见的基本类型
  string hash list set zset
* tcp/ip的粘包问题，怎么解决
* mysql的隔离级别
* 两道算法题 第二道不要求写出来

## 小米 Miot
### 一面
Redis缓存策略
分布式CAP原理
Mysql事务的隔离级别 分别解决了哪些问题？
mysql的存储引擎的区别？
Redis持久化的区别
Grpc的过程
### 二面
怼了项目 凉

* 怎么限制goroutine的数量？
通过channel每次执行的go之前向通道写入值，直到通道满的时候就阻塞了。

* Redis为什么快？
  1. 完全基于内存，绝大部分请求是纯粹的内存操作，非常快速。数据存在内存中，类似于HashMap，HashMap的优势就是查找和操作的时间复杂度都是O(1)；
  2. 采用单线程，避免了不必要的上下文切换和竞争条件，也不存在多进程或者多线程导致的切换而消耗 CPU，不用去考虑各种锁的问题，不存在加锁释放锁操作，没有因为可能出现死锁而导致的性能消耗；
  3. 使用多路I/O复用模型，非阻塞IO；
   多路I/O复用模型是利用 select、poll、epoll 可以同时监察多个流的 I/O 事件的能力，在空闲的时候，会把当前线程阻塞掉，当有一个或多个流有 I/O 事件时，就从阻塞态中唤醒，于是程序就会轮询一遍所有的流（epoll 是只轮询那些真正发出了事件的流），并且只依次顺序的处理就绪的流，这种做法就避免了大量的无用操作。
   多路指多个网络连接 复用指复用同一个线程

## 字节跳动 懂车帝
### 一面 
* https的过程
  1. 客户端请求网址，服务器接收到请求后返回证书公钥
  2. 客户端验证证书的有效性和合法性，然后生成一个随机值
  3. 客户端通过证书的公钥加密随机值，将加密后的密钥发送给服务器
  4. 服务器通过私钥解密密钥，通过解密后的密钥加密要发送的内容
  5. 客户端通过密钥解密接受的内容
* http keep-alive的作用
* 755 是什么权限？ 详细见linux权限入门
 7=rwx=4+2+1 5=rx 
 拥有者可读可写可执行 群组可读可执行 其他组可读可执行 
* linux 查找当前文件后20行 
  tail -n 20 filename
* linux统计字符出现的个数
  grep -o objStr  filename|wc -l
* epoll和poll的区别
  本质都是IO多路复用
* 算法题 10进制转7进制
* 算法题 判断数独的有效性

