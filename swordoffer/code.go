package main

import "fmt"

func main() {
	num := []int{1, 3, 5, 8, 0}
	res := minArray(num)
	fmt.Printf("%v", res)
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
	return 0
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
