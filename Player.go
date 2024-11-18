package main

import "fmt"

func main() {
	var x interface{}

	x = "Hello, World!"

	res, ok := x.(int)

	if ok {
		fmt.Println(res)
	} else {
		fmt.Println("类型断言失败")
	}
}
