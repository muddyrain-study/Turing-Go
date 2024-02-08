package main

import "fmt"

func main() {
	//var arr [3]int
	//fmt.Println(arr[0])
	//fmt.Println(arr[1])
	//fmt.Println(arr[2])

	var arr = [...]int{4, 5, 6}

	arr[0] = 1

	for index, value := range arr {
		fmt.Printf("索引:%d,值:%d \n", index, value)
	}

	a := [2]int{1, 2}
	b := [...]int{1, 2}
	c := [2]int{1, 3}
	fmt.Println(a == b, a == c, b == c) // "true false false"
	//	d := [3]int{1, 2}
	//	fmt.Println(a == d) // 编译错误：无法比较 [2]int == [3]int
}
