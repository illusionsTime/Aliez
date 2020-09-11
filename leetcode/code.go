package main

import (
	"errors"
	"fmt"
)

func main() {
	nums := []int{3, 2, 1, 5, 6, 4}
	HeapSort(nums)
	//mid := size / 2
	fmt.Printf("%v", nums)
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

//leetcode HOT 33
func search(nums []int, target int) int {
	size := len(nums)
	if size == 1 && nums[0] == target {
		return 0
	} else if size == 1 && nums[0] != target {
		return -1
	}
	var mid int
	l, r := 0, size-1
	for l <= r {
		mid = (l + r) / 2
		fmt.Printf("%v", mid)
		if nums[mid] == target {
			return mid
		}
		if nums[0] <= nums[mid] {
			if target < nums[mid] && target >= nums[0] {
				//在前半段查找
				r = mid - 1
			} else {
				l = mid + 1
			}
		} else {
			if target > nums[mid] && target <= nums[size-1] {
				l = mid + 1
			} else {
				r = mid - 1
			}
		}
	}
	return -1
}

//leetcode HOT 139
func wordBreak(s string, wordDict []string) bool {
	size := len(s)
	dp := make([]bool, size+1)
	dp[0] = true
	word := make(map[string]bool, len(wordDict))
	for _, v := range wordDict {
		word[v] = true
	}
	for i := 1; i <= size; i++ {
		for j := 0; j < i; j++ {
			if dp[j] && word[s[j:i]] {
				dp[i] = true
				break
			}
		}
	}
	return dp[size]
}

//leetcode HOT 102
func levelOrder(root *TreeNode) [][]int {
	ret := [][]int{}
	if root == nil {
		return ret
	}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		l := len(queue)
		ans := make([]int, 0)
		for i := 0; i < l; i++ {
			node := queue[i]
			ans = append(ans, node.Val)
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		ret = append(ret, ans)
		queue = queue[l:]
	}
	return ret
}

func levelOrderDFS(root *TreeNode) [][]int {
	ret := [][]int{}
	if root == nil {
		return ret
	}
	var SerchRoot func(root *TreeNode, index int)
	SerchRoot = func(root *TreeNode, index int) {
		if root == nil {
			return
		}
		if len(ret) == index {
			ret = append(ret, []int{})
		}
		ret[index] = append(ret[index], root.Val)
		SerchRoot(root.Left, index+1)
		SerchRoot(root.Right, index+1)
	}
	SerchRoot(root, 0)
	return ret
}

// leetcode 257
func binaryTreePaths(root *TreeNode) []string {
	path := make([]string, 0)
	if root == nil {
		return path
	}
	var DFS func(root *TreeNode, str string)
	DFS = func(root *TreeNode, str string) {
		if root == nil {
			return
		}
		str = fmt.Sprintf("%s%d->", str, root.Val)
		if root.Left == nil && root.Right == nil {
			path = append(path, str[:len(str)-2])
		}
		DFS(root.Left, str)
		DFS(root.Right, str)
	}
	DFS(root, "")
	return path
}

//leetcode 107
func levelOrderBottom(root *TreeNode) [][]int {
	ret := [][]int{}
	if root == nil {
		return ret
	}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		l := len(queue)
		ans := make([]int, 0)
		for i := 0; i < l; i++ {
			node := queue[i]
			ans = append(ans, node.Val)
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		ret = append(ret, ans)
		queue = queue[l:]
	}
	//反转数组
	l := len(ret)
	for i := 0; i < l/2; i++ {
		tmp := ret[i]
		ret[i] = ret[l-i-1]
		ret[l-i-1] = tmp
	}
	return ret
}

//leetcode 144
func preorderTraversal(root *TreeNode) []int {
	ret := make([]int, 0)
	if root == nil {
		return ret
	}
	var PT func(root *TreeNode)
	PT = func(root *TreeNode) {
		ret = append(ret, root.Val)
		if root.Left != nil {
			PT(root.Left)
		}
		if root.Right != nil {
			PT(root.Right)
		}
	}
	PT(root)
	return ret
}

func preorderTraversal2(root *TreeNode) []int {
	ret := make([]int, 0)
	if root == nil {
		return ret
	}
	stack := []*TreeNode{root}
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		ret = append(ret, node.Val)
		if node.Right != nil {
			stack = append(stack, node.Right)
		}
		if node.Left != nil {
			stack = append(stack, node.Left)
		}
	}
	return ret
}

//leetcode 114
func flatten(root *TreeNode) {
	if root == nil {
		return
	}
	stack := []*TreeNode{root}
	var tmp *TreeNode
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if tmp != nil {
			tmp.Right, tmp.Left = nil, node
		}
		left, right := tmp.Left, tmp.Right
		if right != nil {
			stack = append(stack, node.Right)
		}
		if left != nil {
			stack = append(stack, node.Left)
		}
		tmp = node
	}
}

func flatten3(root *TreeNode) {
	if root == nil {
		return
	}
	curr := root
	for curr != nil {
		if curr.Left != nil {
			next := curr.Left
			per := next
			for per.Right != nil {
				per = per.Right
			}
			per.Right = curr.Right
			curr.Left = nil
			curr.Right = next
		}
		curr = curr.Right
	}
}

//leetcode 203
func removeElements(head *ListNode, val int) *ListNode {
	sentry := new(ListNode)
	sentry.Next = head
	p := sentry
	for p.Next != nil {
		if p.Next.Val == val {
			tmp := p.Next
			p.Next = tmp.Next
		} else {
			p = p.Next
		}
	}
	return sentry.Next
}

//leetcode 147
func insertionSortList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	sentry := &ListNode{
		Next: head,
	}
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

//leetcode 160
func getIntersectionNode(headA, headB *ListNode) *ListNode {
	if headA == nil || headB == nil {
		return nil
	}
	l1, l2 := headA, headB
	for l1 != l2 {
		if l1 == nil {
			l1 = headB
		} else {
			l1 = l1.Next
		}

		if l2 == nil {
			l2 = headA
		} else {
			l2 = l2.Next
		}

	}
	return l1
}

//mianshi 2.08
func detectCycle(head *ListNode) *ListNode {
	fast, slow := head, head
	if head == nil {
		return nil
	}
	for fast != nil {
		//避免越界
		if fast.Next == nil || slow.Next == nil {
			return nil
		}
		fast = fast.Next.Next
		slow = slow.Next
		if fast == slow {
			break
		}
	}
	//避免无环时导致越界
	if fast == nil || slow == nil {
		return nil
	}
	slow = head
	for slow != fast {
		slow = slow.Next
		fast = fast.Next
	}
	return fast
}

//leetcode 83
func deleteDuplicates(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	h := head
	var tmp *ListNode
	for h.Next != nil {
		if h.Next.Val == h.Val {
			tmp = h.Next.Next
			h.Next.Next = nil
			h.Next = tmp
		} else {
			h = h.Next
		}
	}
	return head
}

//leetcode 343
func integerBreak(n int) int {
	dp := make([]int, n+1)
	dp[0], dp[1] = 0, 0
	for i := 2; i <= n; i++ {
		max := 0
		for j := 1; j < i; j++ {
			fast := j * (i - j)
			last := j * dp[i-j]
			if fast > last && fast > max {
				max = fast
			} else if last >= fast && last > max {
				max = last
			}
		}
		dp[i] = max
	}
	return dp[n]
}

func solveStep(N int) int {
	if N == 1 {
		return 1
	} else if N == 2 {
		return 1
	} else if N == 3 {
		return 2
	}
	return solveStep(N-1) + solveStep(N-3)
}

//jianzhi 63
func maxProfit(prices []int) int {
	size := len(prices)
	if size == 0 || size == 1 {
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

//jianzhi 47
func maxValue(grid [][]int) int {
	n := len(grid)
	m := len(grid[0])
	dp := make([][]int, n)
	for i := 0; i < m; i++ {
		dp[i] = make([]int, m)
	}
	dp[0][0] = grid[0][0]
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j == 0 && i != 0 {
				dp[i][j] = dp[i-1][j] + grid[i][j]
			} else if i == 0 && j != 0 {
				dp[i][j] = dp[i][j-1] + grid[i][j]
			} else if i != 0 && j != 0 {
				dp[i][j] = max(dp[i][j-1], dp[i-1][j]) + grid[i][j]
			}
		}
	}
	return dp[n][m]
}

//jianzhi 42
func maxSubArray(nums []int) int {
	if len(nums) == 1 {
		return nums[0]
	}
	sum := 0
	x := nums[0]
	for i := 0; i < len(nums); i++ {
		sum = max(sum+nums[i], nums[i])
		x = max(x, sum)
	}
	return x
}

//leetcode 200
func numIslands(grid [][]byte) int {
	n := len(grid)
	if n == 0 {
		return 0
	}
	m := len(grid[0])
	num := 0
	var DFS func(x, y int)
	DFS = func(x, y int) {
		if x < 0 || y < 0 || y >= m || x >= n {
			return
		}
		if grid[x][y] == '0' {
			return
		}
		grid[x][y] = '0'
		DFS(x+1, y)
		DFS(x-1, y)
		DFS(x, y+1)
		DFS(x, y-1)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '1' {
				num += 1
				DFS(i, j)
			}
		}
	}
	return num
}

func judge(root *TreeNode) bool {
	ret := make([]int, 0)
	num := 0
	var jieguo bool
	if root == nil {
		return false
	}
	var PT func(root *TreeNode)
	PT = func(root *TreeNode) {
		if root.Left == nil && root.Right == nil {
			ret = append(ret, root.Val)
			if len(ret) == 2 {
				num = ret[1] - ret[0]
			} else if len(ret) > 2 {
				if root.Val-ret[len(ret)-1] != num {
					jieguo = false
				}
			}
		}
		if root.Left != nil {
			PT(root.Left)
		}
		if root.Right != nil {
			PT(root.Right)
		}
	}
	PT(root)
	if len(ret) <= 2 || jieguo == false {
		return false
	}
	return true
}

//jianzhi 28
func isSymmetric(root *TreeNode) bool {
	if root == nil {
		return true
	}
	return dfsIsSymmetric(root.Left, root.Right)
}
func dfsIsSymmetric(a, b *TreeNode) bool {
	if a == nil && b == nil {
		return true
	}
	if (a == nil && b != nil) || (a != nil && b == nil) || (a.Val != b.Val) {
		return false
	}
	return dfsIsSymmetric(a.Left, b.Right) && dfsIsSymmetric(a.Right, b.Left)
}

//leetcode 34
func searchRange(nums []int, target int) []int {
	if len(nums) == 0 {
		return []int{-1, -1}
	}
	first := findTarget(nums, target, 0)
	last := findTarget(nums, target, 1)
	return []int{first, last}
}
func findTarget(nums []int, target, n int) int { // n区别第一个还是最后一个
	l, r := 0, len(nums)-1
	for l <= r {
		mid := l + (r-l)/2
		if nums[mid] == target {
			if n == 0 && mid > 0 && nums[mid-1] == target { // 第一个等于target
				r = mid - 1
			} else if n == 1 && mid+1 <= r && nums[mid] == nums[mid+1] { // 最后一个等于target
				l = mid + 1
			} else {
				return mid
			}
		} else if nums[mid] > target {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return -1
}

//leetcode 215
func findKthLargest(nums []int, k int) int {
	return 0
}
