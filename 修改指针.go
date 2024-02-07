package main

import (
	"flag"
	"fmt"
)

var mode = flag.String("mode", "", "运行模式可以设置为fast")

func main() {
	//var num = 10
	//modifyFromPoint(num)
	//println("num的值为：", num)
	//
	//var num2 = 22
	//newModifyFromPoint2(&num2)
	//println("num2的值为：", num2)

	//ptr := new(int)
	//*ptr = 100
	//fmt.Println(ptr)
	//fmt.Println(*ptr)

	// 解析命令行参数
	flag.Parse()
	// 输出命令行参数
	fmt.Printf("运行模式为: %s", *mode)
}

func modifyFromPoint(num int) {
	num = 100
	fmt.Println("未使用num的值为：", num)
}
func newModifyFromPoint2(num *int) {
	*num = 100
	fmt.Println("使用num的值为：", *num)
}
