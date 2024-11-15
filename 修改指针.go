package main

import (
	"fmt"
)

func main() {
	//var num = 10
	//modifyFromPoint(num)
	//println("num的值为：", num)
	//
	//newModifyFromPoint2(&num)
	//println("num的值为：", num)

	var ptr *string = new(string)

	*ptr = "吗喽之路GO教程"

	fmt.Printf("%s", *ptr)

}

func modifyFromPoint(num int) {
	num = 100
	fmt.Println("未使用num的值为：", num)
}
func newModifyFromPoint2(num *int) {
	*num = 100
	fmt.Println("使用num的值为：", *num)
}
