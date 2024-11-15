package main

import "fmt"

func main() {
	//var a = [3]int{1, 2, 3}
	//
	//fmt.Println(a, a[:2])

	//var highRiseBuilding [30]int
	//
	//for i := 0; i < len(highRiseBuilding); i++ {
	//	highRiseBuilding[i] = i + 1
	//}
	//
	//fmt.Println(highRiseBuilding[10:15])
	//fmt.Println(highRiseBuilding[20:])
	//fmt.Println(highRiseBuilding[:20])

	//var strList []string
	//var numList []int
	//var numListEmpty = []int{}
	//
	//// 打印三个切片的长度
	//fmt.Println(len(strList), len(numList), len(numListEmpty))
	//// 判定空的结果
	//fmt.Println(strList == nil, numList == nil, numListEmpty == nil)
	//
	//strList = append(strList, "hello")
	//fmt.Println(strList)

	//a := make([]int, 2)
	//b := make([]int, 2, 10)
	//fmt.Println(a, b)
	//a = append(a, 2)
	//b = append(b, 2)
	//fmt.Println(len(a), len(b))

	//var numbers4 = [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	//mySlice := numbers4[4:6]
	//fmt.Println(mySlice)
	//fmt.Println(cap(mySlice))
	//mySlice = mySlice[:cap(mySlice)]
	//fmt.Println(mySlice)
	//fmt.Println(mySlice[3])

	//slice1 := []int{1, 2, 3, 4, 5}
	//slice2 := []int{5, 4, 3}
	//copy(slice2, slice1)
	//fmt.Println(slice2)
	//copy(slice1, slice2)
	//fmt.Println(slice1)

	//var map1 = map[string]int{
	//	"one": 1,
	//	"two": 2,
	//}
	//map1["three"] = 3
	//
	//fmt.Println(map1)

	//map1 := map[string][]int{}
	//map2 := make(map[int][]int)
	//
	//map1["one"] = []int{1, 2, 3}
	//map1["two"] = []int{4, 5, 6}
	//map1["three"] = []int{7, 8, 9}
	//map2[1] = []int{1, 2, 3}
	//
	//fmt.Println(map1, map2)
	//for key, value := range map1 {
	//	fmt.Println(key, value)
	//}
	//
	//delete(map1, "one")
	//fmt.Println(map1)

	var arr []int
	var num *int
	fmt.Printf("%p\n", arr)
	fmt.Printf("%p\n", num)
	arr = append(arr, 1)
	fmt.Printf("%p\n", arr)
}
