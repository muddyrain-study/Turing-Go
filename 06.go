package main

import "fmt"

func length(s string) int {
	println("call length.")
	return len(s)
}

func main() {
	m := map[string]int{
		"one": 1,
		"two": 2,
	}

	for key, value := range m {
		fmt.Println(key, value)
	}
}
