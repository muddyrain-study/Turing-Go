package main

import (
	"flag"
	"fmt"
)

// 定义命令
var mode = flag.String("mode", "", "fast模式能让程序运行的更快")

type NewInt int
type AliasInt = int

func main() {
	var a NewInt
	fmt.Printf("a type: %T\n", a)

	var b AliasInt
	fmt.Printf("b type: %T\n", b)
}
