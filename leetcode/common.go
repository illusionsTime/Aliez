package main

import "errors"

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

//堆 heap
var heap []int

func HeapSort(num []int) {
	var i int
	size := len(num)
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
		child = 2*i + 1
		if child != size-1 && num[child+1] > num[child] {
			child++
		}
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
