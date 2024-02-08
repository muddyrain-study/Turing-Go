package main

import "fmt"

func main() {
	//// 声明一个二维整型数组，两个维度的长度分别是 4 和 2
	//var array [4][2]int
	//// 使用数组字面量来声明并初始化一个二维整型数组
	//array = [4][2]int{{10, 11}, {20, 21}, {30, 31}, {40, 41}}
	//// 声明并初始化数组中索引为 1 和 3 的元素
	//array = [4][2]int{1: {20, 21}, 3: {40, 41}}
	//// 声明并初始化数组中指定的元素
	//array = [4][2]int{1: {0: 20}, 3: {1: 41}}
	//
	//fmt.Println(array[1][0])
	//
	//for index, value := range array {
	//	fmt.Printf("索引:%d,值:%d \n", index, value)
	//}

	// 声明两个二维整型数组 [2]int [2]int
	var array1 [2][2]int
	var array2 [2][2]int
	// 为array2的每个元素赋值
	array2[0][0] = 10
	array2[0][1] = 20
	array2[1][0] = 30
	array2[1][1] = 40
	// 将 array2 的值复制给 array1
	array1 = array2

	fmt.Println(array1)
	fmt.Println(array2)

}
