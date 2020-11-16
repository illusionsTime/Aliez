package main

import "fmt"

/*
张敏阳
*/
func Slove(wide []int, high []int) int {
	size := len(wide)
	if size == 0 {
		return 0
	}
	var sum int
	sum = wide[0] * high[0]
	for i := 1; i < size; i++ {
		right := i + 1
		sumwide := wide[i]
		for high[right] >= high[i] && right < size {
			sumwide += wide[right]
			right++
		}
		left := i - 1
		for high[left] >= high[i] && left >= 0 {
			sumwide += wide[left]
			left--
		}
		if sumwide*high[i] > sum {
			sum = sumwide * high[i]
		}
	}
	return sum
}

func CheckNum(nums []int, target int) []int {
	size := len(nums)
	if size == 0 || nums == nil {
		return []int{0, 0}
	}
	dp := make([]int, size)
	res := make([]int, 2)
	for i := 0; i < size; i++ {
		dp[i] = target - nums[i]
		for j := 0; j < i; j++ {
			if dp[j] == nums[i] {
				res[0] = j
				res[1] = i
			}
		}
	}
	return res
}

func MergeLink(h1 *Node, h2 *Node) *Node {
	sentry1 := h1
	sentry2 := h2
	h := &Node{}
	if sentry1.Val >= sentry2.Val {
		h.Next = sentry1
	} else {
		h.Next = sentry2
	}
	sen := h
	for sentry1.Next != nil && sentry2.Next != nil {
		if sentry1.Val <= sentry2.Val {
			h.Next = sentry1
			sentry1 = sentry1.Next
		} else {
			h.Next = sentry2
			sentry2 = sentry2.Next
		}
	}
	if sentry1.Next == nil && sentry2.Next != nil {
		h.Next = sentry2
	}
	if sentry2.Next == nil && sentry1.Next != nil {
		h.Next = sentry1
	}

	return sen.Next
}

func main() {
	t1 := new(Tree)
	t2 := new(Tree)
	t3 := new(Tree)
	t4 := new(Tree)
	t1.Left = t2
	t1.Right = t3
	t1.Val = 1
	t2.Val = 2
	t2.Left = t4
	t2.Right = nil
	t3.Val = 3
	t3.Left = nil
	t3.Right = nil
	t4.Val = 4
	t4.Left = nil
	t4.Right = nil
	slove(t1)
}

func defer_call() {
	defer func() { fmt.Printf("打印前") }()
	defer func() { fmt.Printf("打印中") }()
	defer func() { fmt.Printf("打印后") }()
	panic("触发异常")
}

func GetUpDay(nums []int) int {
	size := len(nums)
	if size == 0 || size == 1 {
		return 0
	}
	res := 0
	result := 0
	for i := 1; i < size; i++ {
		var num int
		if nums[i] > nums[i-1] {
			res++
			num = res
		} else {
			num = res
			res = 0
		}
		result = max(result, num)
	}
	return result
}

func max(a int, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

func Serch(nums []int) int {
	n := len(nums)
	i := 0
	for i < n {
		if nums[i] != i {
			tmp := nums[i]
			if nums[tmp] == tmp {
				return tmp
			} else {
				nums[i] = nums[tmp]
				nums[tmp] = tmp
			}
		}
	}
	return 0
}

func getmax(num []int) int {
	size := len(num)
	var sum int
	sum = num[0]
	dp := make([]int, size)
	dp[0] = num[0]
	for i := 1; i < size; i++ {
		if dp[i-1]+num[i] > 0 {
			dp[i] = dp[i-1] + num[i]
		}
		sum = max(dp[i], sum)
	}
	return sum
}

type Node struct {
	Next *Node
	Val  int
}

func MergeNode(h1 *Node, h2 *Node, h3 *Node) *Node {
	h := Merge(h1, h2)
	return Merge(h, h3)
}

func Merge(h1 *Node, h2 *Node) *Node {
	sentry := new(Node)
	res := sentry
	for h1 != nil && h2 != nil {
		if h1.Val >= h2.Val {
			sentry.Next = h2
			h2 = h2.Next
		} else {
			sentry.Next = h1
			h1 = h1.Next
		}
		sentry = sentry.Next
	}
	if h1 == nil && h2 != nil {
		sentry.Next = h2
	}
	if h1 != nil && h2 == nil {
		sentry.Next = h1

	}
	return res.Next
}

type Tree struct {
	Val   int
	Left  *Tree
	Right *Tree
}

func slove(t *Tree) {
	if t == nil {
		return
	}
	var DFS func(t *Tree)
	DFS = func(t *Tree) {
		fmt.Printf("%v", t.Val)
		if t.Left != nil {
			DFS(t.Left)
		}
		if t.Right != nil {
			DFS(t.Right)
		}
	}
	DFS(t)
}
