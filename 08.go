package main

import "fmt"

type Color struct {
	R, G, B int
}

type Point struct {
	x int
	y int
}
type Player struct {
	Name        string
	HealthPoint int
	MagicPoint  int
}
type Command struct {
	Name    string
	varRef  *int
	Comment string
}

func main() {
	//color := Color{255, 0, 0}
	//
	//println(color.R, color.G, color.B)
	//
	//point := Point{
	//	x: 1,
	//	y: 2,
	//}
	//
	//println(point.x, point.y)
	//
	//p := new(Player)
	//
	//(*p).Name = "zhang3"
	//(*p).HealthPoint = 100
	//(*p).MagicPoint = 200
	//
	//println(p)
	//var newVersion *int
	//cmd := newCommand("cmd", newVersion, "this is a command")
	//println((cmd))

	msg := struct {
		id   int
		data string
	}{
		1024,
		"hello",
	}
	printMsgType(&msg)
	fmt.Printf("1: %T , %v", msg, msg)
}

func printMsgType(msg *struct {
	id   int
	data string
}) {
	msg.id = 2048
	fmt.Printf("2: %T , %v", msg, msg)
}

func newCommand(name string, varRef *int, comment string) *Command {
	return &Command{
		Name:    name,
		varRef:  varRef,
		Comment: comment,
	}
}
