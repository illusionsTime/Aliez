package main

var a = 0

func main() {
	f := fa()
	g := fa()
	println(f(1))
	println(f(1))
	println(g(1))
	println(g(1))
}

func fa() func(i int) int {
	return func(i int) int {
		println(&a, a)
		a = a + i
		return a
	}
}
