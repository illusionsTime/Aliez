package main

import (
	"errors"
	"fmt"
)

func main() {
	A := []int{1, 12, 3, 43, 23, 28, 17, 9}
	MergeSort(A)
	fmt.Printf("%v", A)
}

//leetcode 面试题03
func findRepeatNumber(nums []int) int {
	n := len(nums)
	m := make(map[int]int, n)
	for k, v := range nums {
		if _, ok := m[v]; ok {
			return v
		}
		m[v] = k
	}
	return 0
}

//leetcode 1111
func maxDepthAfterSplit(seq string) []int {
	ans := new(stack)
	s := make([]int, len(seq))
	for i := 0; i < len(seq); i++ {
		if i == 0 {
			ans.Push(seq[i])
			s[i] = 0
		} else if i > 0 {
			if seq[i] == ans.Top() {
				s[i] = ans.Len() % 2
				ans.Push(seq[i])
			} else {
				ans.Pop()
				s[i] = ans.Len() % 2
			}
		}
	}
	return s
}

//leetcode 289
func gameOfLife(board [][]int) {
	sizex := len(board[0])
	sizey := len(board)
	copydata := make([][]int, sizey)
	for my := 0; my < sizey; my++ {
		copydata[my] = make([]int, sizex)
		for mx := 0; mx < sizex; mx++ {
			copydata[my][mx] = board[my][mx]
		}
	}

	for k := 0; k < sizey; k++ {
		for i := 0; i < sizex; i++ {
			tmp := 0
			//if t <2 t>3 die
			// if t =2 live will live
			// if t=3 all live
			for m := -1; m <= 1; m++ {
				for n := -1; n <= 1; n++ {
					cm := k + m
					cn := i + n
					if !(m == 0 && n == 0) {
						if ((cm >= 0 && cm < sizey) && (cn >= 0 && cn < sizex)) && (copydata[cm][cn] == 1) {
							tmp++
						}
					}

				}
			}
			if (tmp < 2 || tmp > 3) && board[k][i] == 1 {
				board[k][i] = 0
			} else if tmp == 3 && board[k][i] == 0 {
				board[k][i] = 1
			}
		}
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

	//
	for i := 0; i < NumElements; i, REnd = i+1, REnd-1 {
		A[REnd] = TmpArray[REnd]
	}
}
