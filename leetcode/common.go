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
