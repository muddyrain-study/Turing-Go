package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("./main.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var content []byte
	for {
		buf := make([]byte, 128)
		n, err := file.Read(buf[:])
		if err == io.EOF {
			// 读取到文件末尾
			break
		}
		if err != nil {
			fmt.Println("read file err", err)
			return
		}
		content = append(content, buf[:n]...)
	}
	fmt.Println(string(content))
}
