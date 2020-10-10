package main

import "fmt"

func main() {
	nums := []int{1, 2, 3, 4, 5}
	n := exchange(nums)
	fmt.Printf("%v", n)
}

//常用数据结构
/*
linkedlist
tree
*/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

//剑指offer 04
func findNumberIn2DArray(matrix [][]int, target int) bool {
	//以左下角为原点
	i := len(matrix) - 1 //获取右下角y坐标
	j := 0               //获取右下角x坐标
	for i > -1 {
		if j < len(matrix[i]) {
			if target < matrix[i][j] {
				i-- //小于target,向上查找
			} else if target > matrix[i][j] {
				j++ //大于targat,向右查找
			} else if target == matrix[i][j] {
				return true
			}
		} else {
			return false //超出数组返回false
		}
	}
	return false //超出matrix返回false
}

func BinarySearch(array []int, res int, l int, r int) bool {
	mid := (l + r) / 2
	if res < array[0] || res > array[r] {
		return false
	}
	if res == array[mid] || res == array[l] || res == array[r] {
		return true
	}
	if res >= array[l] && res < array[mid] {
		return BinarySearch(array, res, l, mid)
	} else if res <= array[r] && res > array[mid] {
		return BinarySearch(array, res, mid, r)
	} else {
		return false
	}
}

//剑指offer 05
func replaceSpace(s string) string {
	var num []rune
	for _, v := range s {
		if v == ' ' {
			num = append(num, '%', '2', '0')
		} else {
			num = append(num, v)
		}
	}
	return string(num)
}

//剑指offer 07
func buildTree(preorder []int, inorder []int) *TreeNode {
	root := new(TreeNode)
	i := 0
	for ; i < len(inorder); i++ {
		if inorder[i] == preorder[0] {
			root.Val = preorder[0]
			root.Left = buildTree(preorder[1:i+1], inorder[:i])
			root.Right = buildTree(preorder[i+1:], inorder[i+1:])
			return root
		}
	}
	return nil
}

//剑指offer 10
func numWays(n int) int {
	if n == 0 || n == 1 {
		return 1
	}
	dp := make([]int, n+1)
	dp[0] = 1
	dp[1] = 1
	for i := 2; i < n+1; i++ {
		dp[i] = (dp[i-1] + dp[i-2]) % (1e9 + 7)
	}
	return dp[n]
}

//剑指offer 11
func minArray(numbers []int) int {
	l, r := 0, len(numbers)-1
	for l < r {
		mid := l + (r-l)/2
		if numbers[mid] > numbers[r] {
			l = mid + 1
		} else if numbers[mid] <= numbers[r] {
			r = mid
		}
	}
	return numbers[l]
}

//剑指offer 12
func exist(board [][]byte, word string) bool {
	row := len(board)
	line := len(board[0])
	if row == 0 {
		return false
	}

	var DFS func(x, y int, num int) bool
	DFS = func(x, y int, num int) bool {
		if x < 0 || y < 0 || x >= line || y >= row {
			return false
		}
		if board[y][x] == word[len(word)-1] && num == len(word)-1 {
			return true
		}
		tmp := board[y][x]
		board[y][x] = byte(' ')
		if tmp == word[num] {
			num++
			if DFS(x+1, y, num) || DFS(x-1, y, num) || DFS(x, y+1, num) || DFS(x, y-1, num) {
				return true
			}
		}
		board[y][x] = tmp
		return false
	}
	for i := 0; i < row; i++ {
		for j := 0; j < line; j++ {
			if DFS(j, i, 0) {
				return true
			}
		}
	}
	return false
}

//剑指offer 13
func movingCount(m int, n int, k int) int {
	if k == 0 {
		return 1
	}
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	var sum int
	var DFS func(x int, y int)
	DFS = func(x int, y int) {
		if x < 0 || y < 0 || x >= m || y >= n {
			return
		}
		if dp[x][y] == 1 {
			return
		}
		if Digitstogether(x)+Digitstogether(y) > k {
			return
		}
		dp[x][y] = 1
		sum++
		DFS(x, y+1)
		DFS(x, y-1)
		DFS(x+1, y)
		DFS(x-1, y)
	}
	DFS(0, 0)
	return sum
}

//位数相加
func Digitstogether(num int) int {
	sum := 0
	for num > 0 {
		sum += num % 10
		num = num / 10
	}
	return sum
}

//剑指offer 15
func hammingWeight(num uint32) int {
	count := 0
	for 0 < num {
		if num%2 == 1 {
			count++
		}
		num /= 2
	}
	return count
}

//剑指offer 15
func printNumbers(n int) []int {
	s := 1
	for i := 1; i <= n; i++ {
		s = 10 * s
	}
	s = s - 1
	num := make([]int, s)
	for i := 0; i < s; i++ {
		num[i] = i + 1
	}
	return num
}

//剑指offer 21
func exchange(nums []int) []int {
	size := len(nums)
	if size == 1 || nums == nil {
		return nums
	}
	r, l := size-1, 0
	for l < r {
		if nums[l]%2 == 0 && nums[r]%2 != 0 {
			var tmp int
			tmp = nums[l]
			nums[l] = nums[r]
			nums[r] = tmp
		}
		if nums[l]%2 != 0 {
			l++
		}
		if nums[r]%2 == 0 {
			r--
		}
		fmt.Printf("%v", nums)
	}
	return nums
}

//剑指offer 29
func spiralOrder(matrix [][]int) []int {
	if matrix == nil || len(matrix) == 0 || len(matrix) == 0 {
		return []int{}
	}
	top := 0
	hsize := len(matrix)
	lsize := len(matrix[0])
	left := 0

	bottom := hsize - 1
	right := lsize - 1
	index := 0
	x, y := 0, 0
	sum := make([]int, hsize*lsize)
	for bottom >= top && right >= left {
		for x = left; x <= right; x++ {
			sum[index] = matrix[top][x]
			index++
		}
		for y = top + 1; y <= bottom; y++ {
			sum[index] = matrix[y][right]
			index++
		}
		if bottom > top && right > left {
			for x = right - 1; x > left; x-- {
				sum[index] = matrix[bottom][x]
				index++
			}
			for y = bottom; y > top; y-- {
				sum[index] = matrix[y][left]
				index++
			}
		}
		left++
		right--
		top++
		bottom--
	}
	return sum
}

//剑指offer 26
func isSubStructure(A *TreeNode, B *TreeNode) bool {
	if A == nil || B == nil {
		return false
	}
	var judge func(A *TreeNode, B *TreeNode) bool
	judge = func(A *TreeNode, B *TreeNode) bool {
		if B == nil {
			return true
		}
		if A == nil && B != nil {
			return false
		}
		if A.Val != B.Val {
			return false
		}
		return judge(A.Left, B.Left) && judge(A.Right, B.Right)
	}
	return judge(A, B) || isSubStructure(A.Left, B) || isSubStructure(A.Right, B)
}

//剑指offer 27
func mirrorTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	root.Left, root.Right = root.Right, root.Left
	mirrorTree(root.Left)
	mirrorTree(root.Right)
	return root
}

//剑指offer 31
func validateStackSequences(pushed []int, popped []int) bool {
	stack := make([]int, 0, len(pushed))
	i := 0
	for _, v := range pushed {
		stack = append(stack, v)
		for len(stack) != 0 && stack[len(stack)-1] == popped[i] {
			stack = stack[:len(stack)-1]
			i++
		}
	}
	if len(stack) == 0 {
		return true
	} else {
		return false
	}

}

//剑指offer 32
func levelOrder(root *TreeNode) []int {
	sum := make([]int, 0)
	if root == nil {
		return sum
	}
	var queue []*TreeNode
	queue = append(queue, root)
	for len(queue) != 0 {
		l := len(queue)
		for i := 0; i < l; i++ {
			node := queue[i]
			sum = append(sum, node.Val)
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		queue = queue[l:]
	}
	return sum
}

//剑指offer 33 二叉搜索树的后序遍历
func verifyPostorder(postorder []int) bool {
	size := len(postorder)
	var slove func(num []int, i, j int) bool
	slove = func(num []int, i, j int) bool {
		if i >= j {
			return true
		}
		piovt := num[j]
		k := i
		for num[k] < piovt {
			k++
		}
		m := k
		for num[k] > piovt {
			k++
		}
		return k == j && slove(num, i, m-1) && slove(num, m, j-1)
	}
	return slove(postorder, 0, size-1)
}

//剑指offer 40
func getLeastNumbers(arr []int, k int) []int {
	if k >= len(arr) {
		return arr
	}
	return quickselect(arr, 0, len(arr)-1, k)
}

func quickselect(arr []int, left, right int, k int) []int {
	if left < right {
		index := partition(arr, left, right)
		if index == k {
			return arr[:k]
		} else if index > k {
			return quickselect(arr, left, index-1, k)
		} else if index < k {
			return quickselect(arr, index+1, right, k)
		}
	}
	return arr[:k]
}

func partition(num []int, l, r int) int {
	poivt := num[l]
	index := l + 1
	for i := index; i <= r; i++ {
		if num[i] < poivt {
			num[i], num[index] = num[index], num[i]
			index++
		}
	}
	num[l], num[index-1] = num[index-1], num[l]
	return index - 1
}

//剑指offer 57
func twoSum(nums []int, target int) []int {
	if len(nums) < 2 {
		return nil
	}
	i, j := 0, len(nums)-1
	for i != j {
		if nums[i]+nums[j] == target {
			return []int{nums[i], nums[j]}
		} else if nums[i]+nums[j] < target {
			i++
		} else if nums[i]+nums[j] > target {
			j--
		}
	}
	return nil
}
