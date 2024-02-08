package main

import "fmt"

func main() {
	//var a = [3]int{1, 2, 3}
	//fmt.Println(a, a[1:2])
	//var highRiseBuilding [30]int
	//for i := 0; i < 30; i++ {
	//	highRiseBuilding[i] = i + 1
	//}
	//// 区间
	//fmt.Println(highRiseBuilding[10:15])
	//// 中间到尾部的所有元素
	//fmt.Println(highRiseBuilding[20:])
	//// 开头到中间指定位置的所有元素
	//fmt.Println(highRiseBuilding[:2])
	//highRiseBuilding[10:20][2] = 1111111

	// 声明字符串切片
	var strList []string
	// 声明整型切片
	var numList []int
	// 声明一个空切片
	var numListEmpty = []int{}
	// 输出3个切片
	fmt.Println(strList, numList, numListEmpty)
	// 输出3个切片大小
	fmt.Println(len(strList), len(numList), len(numListEmpty))
	// 切片判定空的结果
	fmt.Println(strList == nil)
	fmt.Println(numList == nil)
	fmt.Println(numListEmpty == nil)

	strList = append(strList, "哈哈哈")
	fmt.Println(strList)

}
