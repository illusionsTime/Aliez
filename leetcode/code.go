package main

import (
	"errors"
	"fmt"
)

func main() {
	A := []int{1, 12, 3, 43, 23, 28, 17, 9}
	InsertSort(A)
	fmt.Printf("%v", A)
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

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */

type ListNode struct {
	Val  int
	Next *ListNode
}

func mergeKLists(lists []*ListNode) *ListNode {
	r := len(lists)
	return msortList(lists, 0, r-1)
}

func msortList(lists []*ListNode, l int, r int) *ListNode {
	if l > r {
		return nil
	}
	if l == r {
		return lists[l]
	}
	mid := (l + r) >> 1
	return mergeTwoLists(msortList(lists, l, mid), msortList(lists, mid+1, r))
}

func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
	prevhead := new(ListNode)
	prevhead.Val = -1
	prev := prevhead
	for l1 != nil && l2 != nil {
		if l1.Val <= l2.Val {
			prev.Next = l1
			l1 = l1.Next
		} else {
			prev.Next = l2
			l2 = l2.Next
		}
		prev = prev.Next
	}

	if l1 == nil && l2 != nil {
		prev.Next = l2
	} else {
		prev.Next = l1
	}
	return prevhead.Next
}
