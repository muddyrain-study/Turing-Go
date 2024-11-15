package main

import (
	"fmt"
	"strconv"
)

func main() {
	newStr := "1"

	intValue, _ := strconv.Atoi(newStr)

	fmt.Println(intValue)

	newInt := 1

	strValue := strconv.Itoa(newInt)
	fmt.Printf("%T = %v \n", strValue, strValue)

	// string to float
	newStr2 := "1.234"
	parseFloatValue, _ := strconv.ParseFloat(newStr2, 64)
	fmt.Printf("%T = %v \n", parseFloatValue, parseFloatValue)
	// float to string
	floatValue := 1.234
	parseStringValue := strconv.FormatFloat(floatValue, 'f', -1, 64)
	fmt.Printf("%T = %v \n", parseStringValue, parseStringValue)

}
