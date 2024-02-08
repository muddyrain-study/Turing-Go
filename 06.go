package main

import "fmt"

func main() {
	var numbers4 = [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	myslice := numbers4[4:6]
	//这打印出来长度为2
	fmt.Printf("myslice为 %d, 其长度为: %d\n", myslice, len(myslice))
	myslice = myslice[:cap(myslice)]
	//为什么 myslice 的长度为2，却能访问到第四个元素
	fmt.Printf("myslice的第四个元素为: %d", myslice[3])
}
