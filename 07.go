package main

import "fmt"

func main() {
	fmt.Println(max(10, 20))
	fmt.Println(max(-1, -2))

	fmt.Println(test(1, 2, "求和"))
}
func test(x, y int, s string) (int, string) {
	n := x + y
	return n, fmt.Sprintf("%s：%d", s, n)
}
func max(n1, n2 int) int {
	if n1 > n2 {
		return n1
	}
	return n2
}
