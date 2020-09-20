package main

import (
	"errors"
)

//栈
type stack []interface{}

func (s stack) IsEmpty() bool {
	return len(s) == 0
}

func (s *stack) Pop() (interface{}, error) {
	theStack := *s
	if len(theStack) == 0 {
		return 0, errors.New("Out of index, len is 0")
	}
	value := theStack[len(theStack)-1]
	*s = theStack[:len(theStack)-1]
	return value, nil
}

func (s *stack) Push(val interface{}) {
	*s = append(*s, val)
}

func (s stack) Value(i int) interface{} {
	return s[i]
}

func (s stack) Top() interface{} {
	if len(s) == 0 {
		return nil
	}
	return s[len(s)-1]
}

func (s *stack) Len() int {
	return len(*s)
}

//二叉树
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

//链表
type ListNode struct {
	Val  int
	Next *ListNode
}

func max(a int, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

func min(a int, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}

//堆排序
var heap []int

func HeapSort(num []int) {
	var i int
	size := len(num)
	//从最深处父节点构造堆
	for i = size / 2; i >= 0; i-- {
		PercDown(num, i, size)
	}
	for i = size - 1; i > 0; i-- {
		Swap(num[0], num[i])
		PercDown(num, 0, i)
	}
}

func PercDown(num []int, i int, size int) {
	var child int
	var tmp int
	for tmp = num[i]; 2*i+1 < size; i = child {
		child = 2*i + 1 //child是左儿子
		//找到更大的儿子节点
		if child != size-1 && num[child+1] > num[child] {
			child++
		}
		//如果当前父节点小于儿子节点，交换位置
		if tmp < num[child] {
			num[i] = num[child]
		} else {
			break
		}
	}
	num[i] = tmp
}

func Swap(a interface{}, b interface{}) {
	var tmp interface{}
	tmp = a
	a = b
	b = tmp
}

//插入排序
func InsertSort(A []int) {
	var tmp, j int
	N := len(A)
	for i := 0; i < N; i++ {
		tmp = A[i]
		for j = i; j > 0 && A[j-1] > tmp; j-- {
			A[j] = A[j-1]
		}
		A[j] = tmp
	}
}

//归并排序
func MergeSort(A []int) {
	//分配额外空间
	N := len(A)
	tmp := make([]int, N)
	if tmp != nil {
		MSort(A, tmp, 0, N-1)
	} else {
		errors.New("No space for tmp array!")
	}

}

func MSort(A []int, tmp []int, Left int, Right int) {
	var Center int
	if Left < Right {
		//拆分左右数组 递归进行排序
		Center = (Left + Right) / 2
		MSort(A, tmp, Left, Center)
		MSort(A, tmp, Center+1, Right)
		//合并左右数组
		Merge(A, tmp, Left, Center+1, Right)
	}
}

func Merge(A []int, TmpArray []int, Lptr int, Rptr int, REnd int) {
	LEnd := Rptr - 1
	TmpPos := Lptr
	NumElements := REnd - Lptr + 1

	//
	for Lptr <= LEnd && Rptr <= REnd {
		if A[Lptr] <= A[Rptr] {
			TmpArray[TmpPos] = A[Lptr]
			TmpPos++
			Lptr++
		} else {
			TmpArray[TmpPos] = A[Rptr]
			TmpPos++
			Rptr++
		}
	}
	//
	for Lptr <= LEnd {
		TmpArray[TmpPos] = A[Lptr]
		TmpPos++
		Lptr++
	}
	//
	for Rptr <= REnd {
		TmpArray[TmpPos] = A[Rptr]
		TmpPos++
		Rptr++
	}

	//重新填入A数组
	for i := 0; i < NumElements; i, REnd = i+1, REnd-1 {
		A[REnd] = TmpArray[REnd]
	}
}

//希尔排序
func ShellSort(num []int) {
	var tmp int
	var i, j int
	N := len(num)
	for Increment := N / 2; Increment > 0; Increment = Increment / 2 {
		for i = Increment; i < N; i++ {
			tmp = num[i]
			for j = i; j >= Increment; j -= Increment {
				if tmp < num[j-Increment] {
					num[j] = num[j-Increment]
				} else {
					break
				}
			}
			num[j] = tmp
		}
	}
}

//快速排序
func QSort(num []int, left, right int) {
	if right > left {
		pivot := partition(num, left, right)
		QSort(num, left, pivot-1)
		QSort(num, pivot+1, right)
	}
}

func partition(list []int, low, high int) int {
	pivot := list[low] //导致 low 位置值为空
	for low < high {
		//high指针值 >= pivot high指针👈移
		for low < high && pivot <= list[high] {
			high--
		}
		//填补low位置空值
		//high指针值 < pivot high值 移到low位置
		//high 位置值空
		list[low] = list[high]
		//low指针值 <= pivot low指针👉移
		for low < high && pivot >= list[low] {
			low++
		}
		//填补high位置空值
		//low指针值 > pivot low值 移到high位置
		//low位置值空
		list[high] = list[low]
	}
	//pivot 填补 low位置的空值
	list[low] = pivot
	return low
}
