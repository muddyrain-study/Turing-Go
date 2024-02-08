package main

import (
	"bufio"
	"fmt"
	"os"
)

var level = 1
var ex = 0

func main() {
	fmt.Println("请输入你的角色名字")
	//捕获标准输入，并转换为字符串
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	//删除最后的\n
	name := input[:len(input)-1]
	fmt.Printf("角色创建成功,%s,欢迎你来到张三游戏,目前角色等级%d \n", name, level)
	s := `你遇到了一个怪物，请选择是战斗还是逃跑?
	1.战斗
	2.逃跑`
	fmt.Printf("%s \n", s)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		selector := input[:len(input)-1]
		switch selector {
		case "1":
			ex += 10
			fmt.Printf("杀死了怪物，获得了%d经验 \n", ex)
			computeLevel()
			fmt.Printf("您现在的等级为%d \n", level)
		case "2":
			fmt.Printf("你选择了逃跑\n")
			fmt.Printf("%s \n", s)
		case "exit":
			fmt.Println("你退出了游戏")
			//退出
			os.Exit(1)
		default:
			fmt.Println("你的输入我不认识，请重新输入")
		}
	}
}

func computeLevel() {
	if ex < 20 {
		level = 1
	} else if ex < 40 {
		level = 2
	} else if ex < 200 {
		level = 3
	} else {
		level = 4
	}
}
