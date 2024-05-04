package main

import (
	"fmt"
	"math"
)

func main() {
	res := func(x float64) float64 {
		sqrt := math.Sqrt(x)
		fmt.Println(sqrt)
		return sqrt
	}(9)
	fmt.Println(res)
}
