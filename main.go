package main

import (
	"bytes"
	"fmt"
)

func main() {
	data := "123456789"
	//通过[]byte创建Reader
	re := bytes.NewReader([]byte(data))

	buf := make([]byte, 2)

	re.Seek(0, 0)
	//设置偏移量
	for {
		//一个字节一个字节的读
		b, err := re.ReadByte()
		if err != nil {
			break
		}
		fmt.Println(string(b))
	}
	fmt.Println("----------------")

	re.Seek(0, 0)
	off := int64(0)
	for {
		//指定偏移量读取
		n, err := re.ReadAt(buf, off)
		if err != nil {
			break
		}
		off += int64(n)
		fmt.Println(off, string(buf[:n]))
	}

}
