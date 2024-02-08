package main

import (
	"fmt"
	"strconv"
)

func main() {
	newStr1 := "1"
	intValue, _ := strconv.Atoi(newStr1)
	println(intValue)

	newStr2 := 2

	intValue2 := strconv.Itoa(newStr2)
	fmt.Printf("%T %v \n", intValue2, intValue2)

	// string 转 float
	newStr3 := "3.1415926"
	floatValue, _ := strconv.ParseFloat(newStr3, 64)
	fmt.Printf("%T %v \n", floatValue, floatValue)

	// float 转 string
	newStr4 := 3.1415926
	floatValue2 := strconv.FormatFloat(newStr4, 'f', 2, 64)
	fmt.Printf("%T %v \n", floatValue2, floatValue2)

}
