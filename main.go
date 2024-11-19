package main

import "os"

func write() {
	err := os.WriteFile("./yyy.txt", []byte("你好，世界"), 0666)
	if err != nil {
		panic(err)
	}
}
func read() {
	content, err := os.ReadFile("./yyy.txt")
	if err != nil {
		panic(err)
	}
	println(string(content))
}
func main() {
	write()
	read()
}
