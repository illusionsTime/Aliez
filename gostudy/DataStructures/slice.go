package main

import "fmt"

func main() {
	a := make([]int, 20)
	fmt.Println(len(a), cap(a))
	b := make([]int, 42)
	a = append(a, b...)
	fmt.Println(len(a), cap(a))
}
