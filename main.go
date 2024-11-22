package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	var (
		name    string
		age     int
		married bool
	)
	r := strings.NewReader("10 false 张三")

	n, err := fmt.Fscanf(r, "%d %t %s", &age, &married, &name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fscanf:%v\n", err)
	}

	fmt.Println(name, age, married)

	fmt.Println(n)
}
