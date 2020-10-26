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

type Node struct {
	Next *Node
	Val  int
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
	defer_call()
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
