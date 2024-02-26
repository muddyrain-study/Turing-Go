package main

func length(s string) int {
	println("call length.")
	return len(s)
}

func main() {
	s := "abcd"
	// 这样写会多次调佣length函数
	for i, n := 0, length(s); i < n; i++ {
		println(i, s[i])
	}
}
