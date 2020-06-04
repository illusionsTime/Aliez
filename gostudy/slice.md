## 了解slice

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

#### make和var创造的区别

前者的切片指针有分配，后者的内部指针为0
```
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
```
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

```
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