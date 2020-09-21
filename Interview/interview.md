<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [滴滴](#滴滴)
  - [一面](#一面)
  - [二面](#二面)
  - [三面](#三面)
- [富途](#富途)
  - [一面](#一面-1)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 滴滴
### 一面
主要是项目以及部分基础知识，算法的话是股票问题，还是很简单的动态规划问题  
剑指offer.63

假设把某股票的价格按照时间先后顺序存储在数组中，请问买卖该股票一次可能获得的最大利润是多少？

示例 1:
```
输入: [7,1,5,3,6,4]
输出: 5
```
解释: 在第 2 天（股票价格 = 1）的时候买入，在第 5 天（股票价格 = 6）的时候卖出，最大利润 = 6-1 = 5 。  
     注意利润不能是 7-1 = 6, 因为卖出价格需要大于买入价格。  

思路：在第i天时，有两种选择,不出售股票和出售股票，那么不出售股票的利润就是前一天的利润dp[i-1],  
出售的话期利润就是num[i]-num[j],在第j天购买股票，如何判断j是哪一天我们通过一个常量来记录前面遍历过的最小值。

可以写出状态转移方程dp[i]=max(dp[i-1],num[i]-num[j])

```go
func maxProfit(prices []int) int {
	size := len(prices)
    if size==0||size==1{
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
```
* sql题 查找当前班级所有科目平均分大于90的学生姓名  
``` sql 
select name from table group by name having avg(score)>90;
```

### 二面 
二面主要问了项目以及一道算法题和场景设计题，我个人觉得不具有总结性

### 三面
三面的话感觉状态太差，一道送分题变成送命题，后面面试官也就没问什么就结束了

剑指offer 29 顺时针打印矩阵
这道题应该采用暴力的按层模拟打印即可，需要注意的是打印的边界。

![](../views/顺时针打印矩阵.png)

```go
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
```

## 富途
### 一面 
主要是三道算法题 第一二道要写出来，第三道讲思路
先讲第三道吧  
设计一个算法，找出数组中最小的k个数。以任意顺序返回这k个数均可，时间复杂度O(nlogk) 

个人首先想到的是堆排 构建一个小顶堆，每次将堆顶元素滞后，重新下滤，这样得到的是排序后的k个数

第二种方法是快排的变种，找出枢纽元的index与k的关系，相等则返回，index>k证明在枢纽元前边，反之则在后边，时间复杂度应该是等比数列的求和O(n),但是存在最坏情况O(n^2),对于k较大时，这应该是最高效的解法。

```go
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

```

第三种方法是如果数组范围极为有限，通过一个桶即可实现一次遍历完成。

第二道算法题 leetcode209
给定一个含有 n 个正整数的数组和一个正整数 s ，找出该数组中满足其和 ≥ s 的长度最小的 连续 子数组，并返回其长度。如果不存在符合条件的子数组，返回 0。

示例：
```
输入：s = 7, nums = [2,3,1,2,4,3]
输出：2
```
解释：子数组 [4,3] 是该条件下的长度最小的子数组。

思路 暴力法需要O(n^2)，想更进一步可以用二分查找，最好的方法是双指针法时间复杂度为O(n)

```go
func minSubArrayLen(s int, nums []int) int {
	size := len(nums)
	if size == 0 {
		return 0
	}
	sum := 0
	ans := size + 1
	start, end := 0, 0
	for end < size {
		sum += nums[end]
		for sum >= s {
			ans = min(ans, end-start+1)
			sum -= nums[start]
			start++
		}
		end++
	}
	if ans == size+1 {
		return 0
	}
	return ans
}
```

