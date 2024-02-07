package main

import "fmt"

func main() {
	//var cat int = 1
	//
	//var str string = "banana"
	//ptr := &cat
	//
	//fmt.Printf("%p %p \n ", &cat, &str)
	//fmt.Println(*ptr)

	var room = 10
	var ptr = &room

	fmt.Printf("%p \n", ptr)
	fmt.Printf("%T ,%p\n", ptr, ptr)

	fmt.Println("room 变量的地址是：", ptr)
	fmt.Println("room 变量的值是：", *ptr)

}
