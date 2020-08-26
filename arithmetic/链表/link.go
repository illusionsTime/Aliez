package main

import "fmt"

func main() {
	d := &ListNode{
		Val:  3,
		Next: nil,
	}
	c := &ListNode{
		Val:  1,
		Next: d,
	}

	b := &ListNode{
		Val:  2,
		Next: c,
	}

	a := &ListNode{
		Val:  4,
		Next: b,
	}
	fmt.Scanf("%v", &a)
	n := insertionSortList(a)
	for n != nil {
		fmt.Printf("%v", n)
		n = n.Next
	}
}

type ListNode struct {
	Val  int
	Next *ListNode
}

func insertionSortList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	sentry := new(ListNode)
	sentry.Next = head
	p := head.Next
	head.Next = nil
	for p != nil {
		prev := sentry
		next := p.Next
		for prev.Next != nil && prev.Next.Val <= p.Val {
			prev = prev.Next
		}
		// prev.Next.Val > p.Val
		p.Next = prev.Next
		prev.Next = p
		p = next
	}
	return sentry.Next
}
