package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	err := firstCheckError()
	if err != nil {
		fmt.Println(err)
		goto onExit
	}
	err = secondCheckError()
	if err != nil {
		fmt.Println(err)
		goto onExit
	}
	fmt.Println("done")
	return
onExit:
	exitProcess()
}

func secondCheckError() interface{} {
	return errors.New("错误2")
}

func exitProcess() {
	fmt.Println("exit")
	//退出
	os.Exit(1)
}

func firstCheckError() interface{} {
	return errors.New("错误1")
}
