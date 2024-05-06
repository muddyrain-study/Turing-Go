package main

import "fmt"

type Point struct {
	x int
	y int
}
type Color struct {
	r, g, b int
}
type Player struct {
	Name        string
	HealthPoint int
	MagicPoint  int
}

func main() {
	var p Point
	p.x = 10
	p.y = 20
	fmt.Printf("p = %v, x=%d,y=%d \n", p, p.x, p.y)

	var p1 = Point{
		x: 10,
		y: 20,
	}
	fmt.Printf("p1 = %v, x=%d,y=%d \n", p1, p1.x, p1.y)

	tank := new(Player)
	(*tank).Name = "Canon"
	tank.HealthPoint = 300
	tank.MagicPoint = 100

	fmt.Println(tank)
}
