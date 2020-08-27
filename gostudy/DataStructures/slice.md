## 了解slice

### 数组
数组是由相同类型元素的集合组成的数据结构，计算机会为数组分配一块连续的内存来保存其中的元素

对于这个连续内存，一个由字面量组成的数组，根据数组元素数量的不同，编译器会在负责初始化字面量的 cmd/compile/internal/gc.anylit 函数中做两种不同的优化：

1. 当元素数量小于或者等于 4 个时，会直接将数组中的元素放置在栈上；
2. 当元素数量大于 4 个时，会将数组中的元素放置到静态区并在运行时取出（copy到栈上）

[详情可以参考这里](https://github.com/golang/go/blob/f07059d949057f414dd0f8303f93ca727d716c62/src/cmd/compile/internal/gc/sinit.go#L875-L967)
```go

```

### 切片

#### 底层结构

源码src/runtime/slice.go

```
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}
```
切片是由一个指向底层数组的指针，切片的当前大小，底层数组大小组成。当len增长超过cap时，会申请一个更大容量的底层数组，并将数据从老数组复制到新申请的数组中。


#### 扩容
我们使用append向切片内追加数据，当切片容量不足时，就会调用growslice对切片进行扩容。
golang源码里面写了*growslice handles slice growth during append. It is passed the slice element type, the old slice, and the desired new minimum capacity, and it returns a new slice with at least that capacity, with the old data copied into it.*

growslice处理附加期间的切片生长。它传递了片元素类型、旧片和所需的新最小容量，并返回一个至少具有该容量的新片，其中复制了旧数据。
```go
func growslice(et *_type, old slice, cap int) slice {
	...
	newcap := old.cap
	doublecap := newcap + newcap
	if cap > doublecap {
		newcap = cap
	} else {
		if old.len < 1024 {
			newcap = doublecap
		} else {
			// Check 0 < newcap to detect overflow
			// and prevent an infinite loop.
			for 0 < newcap && newcap < cap {
				newcap += newcap / 4
			}
			// Set newcap to the requested cap when
			// the newcap calculation overflowed.
			if newcap <= 0 {
				newcap = cap
			}
		}
	}
	...
}
```
在分配内存空间之前需要先确定新的切片容量，Go 语言根据切片的当前容量选择不同的策略进行扩容：

* 如果期望容量大于当前容量的两倍就会使用期望容量；
* 如果当前切片的长度小于 1024 就会将容量翻倍；
* 如果当前切片的长度大于 1024 就会每次增加 25% 的容量，直到新容量大于期望容量；

```go
func main() {
	a := make([]int, 20, 40)
	fmt.Println(len(a), cap(a))
	b := make([]int, 42)
	a = append(a, b...)
	fmt.Println(len(a), cap(a))
}
```
如上，按我们的理解得到的结果应该是62，80，实际上也确实是这样。但是
```go
func main() {
	a := make([]int, 20)
	fmt.Println(len(a), cap(a))
	b := make([]int, 42)
	a = append(a, b...)
	fmt.Println(len(a), cap(a))
}
```
如上，得到的结果却是62，64
这是因为内存分配的问题，预估容量>2×20,确实是62，但是go中int是8个字节，8×62=496，会选择512的内存空间，所以512 / 8 = 64，即cap(a) = 64
#### make和var创造的区别

前者的切片指针有分配，后者的内部指针为0
```go
var a []int
b:=make([]int,0)
if a==nil{
    ok!
}
if b==nil{
    false!
}//b的底层数组大小为0，但是不为nil
```
也就是说var创造的切片是nil切片（底层数组未分配，指针指向nil）
```
+-------------+    
|pointer=nil  |
+-------------+
|len=0        |
+-------------+
|cap=0        |
+-------------+
```
而make创造的是空切片(底层数组为空，底层数组指针非空)
```
+-------------+         +---------+
|pointer=nil  |---------|  Array  |
+-------------+         +---------+
|len=0        |
+-------------+
|cap=0        |
+-------------+
```

可以看下makeslice的源码
```go
func makeslice(et *_type, len, cap int) unsafe.Pointer {
	mem, overflow := math.MulUintptr(et.size, uintptr(cap))
	if overflow || mem > maxAlloc || len < 0 || len > cap {
		// NOTE: Produce a 'len out of range' error instead of a
		// 'cap out of range' error when someone does make([]T, bignumber).
		// 'cap out of range' is true too, but since the cap is only being
		// supplied implicitly, saying len is clearer.
		// See golang.org/issue/4085.
		mem, overflow := math.MulUintptr(et.size, uintptr(len))
		if overflow || mem > maxAlloc || len < 0 {
			panicmakeslicelen()
		}
		panicmakeslicecap()
	}

	return mallocgc(mem, et, true)
}
```

在runtime/malloc.go中

```go
func mallocgc(size uintptr, typ *_type, needzero bool) unsafe.Pointer {
    ...
    if size == 0 {
		return unsafe.Pointer(&zerobase)
	}
    ...
```
可以看到，当size为0，指向固定的zerobase的全局变量的地址

#### append扩展切片引发的问题（二义性）

多个切片共享底层数组,其中一个切片的append操作可能引发以下两种情况
```
a:=[]int{0,1,2,3,4,5,6,7}
b:=a[0:4]

```

* append追加的元素没有超过底层数组的容量，这是append操作会直接操作共享的底层数组。如果其他切片有引用数组被覆盖的元素，则会导致其他切片的值也隐式的发生变化

* append追加的元素和原来的元素加起来超过了底层数组的容量，则会重新申请新数组，并将原来数组的值拷贝到新数组。


## make和new的区别

make 的作用是初始化内置的数据结构，也就是我们在前面提到的切片、哈希表和 Channel2；
new 的作用是根据传入的类型分配一片内存空间并返回指向这片内存空间的指针3；