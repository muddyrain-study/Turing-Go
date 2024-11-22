package log

import (
	"fmt"
	"log"
	"testing"
)

func TestLog(t *testing.T) {
	log.Print("this is a log")
	log.Printf("this is a log: %d", 100) // 格式化输出
	name := "john"
	age := 20
	log.Println(name, " ", age)
}

func TestLogPanic(t *testing.T) {
	defer fmt.Println("发生了 panic错误！")
	log.Print("this is a log")
	log.Panic("this is a panic log ")
	fmt.Println("运行结束。。。")
}

func TestLogFatal(t *testing.T) {
	defer fmt.Println("defer。。。")
	log.Print("this is a log")
	log.Fatal("this is a fatal log")
	fmt.Println("运行结束。。。")
}

func TestLogFlags(t *testing.T) {
	//i := log.Flags()
	//fmt.Printf("i: %v\n", i)
	//log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	//log.Print("this  is a log")

	s := log.Prefix()
	fmt.Printf("s: %v\n", s)
	log.SetPrefix("[MyLog] ")
	s = log.Prefix()
	fmt.Printf("s: %v\n", s)
	log.Print("this is a log...")
}
