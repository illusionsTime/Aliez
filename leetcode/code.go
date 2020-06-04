package main

func main() {

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

func singleNumbers(nums []int) []int {
	a:=0
    for _,v:=range nums{
        a^=v
    }
   
}