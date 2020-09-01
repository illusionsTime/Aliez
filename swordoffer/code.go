package main

func main() {

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

}
