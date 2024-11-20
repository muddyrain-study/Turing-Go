package main

import (
	"bufio"
	"fmt"
	"net"
)

var wrapperPath = "/Users/qiushunmeng/Movies/GoLang/课程资料/"

func process(conn net.Conn) {
	defer conn.Close() // 关闭连接
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}
		receiveStr := string(buf[:n])
		fmt.Println("收到client端发来的数据：", receiveStr)
		// 读取文件内容
		conn.Write([]byte(receiveStr))
	}
}
func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process(conn)

	}

}
