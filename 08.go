package main

import "fmt"

type Bag struct {
	items []int
}

func (b *Bag) Add(item int) {
	b.items = append(b.items, item)
}

type Property struct {
	Value int
}

func (p *Property) GetValue() int {
	return p.Value
}
func (p *Property) SetValue(v int) {
	p.Value = v
}

type Point struct {
	x int
	y int
}

func (p Point) Add(other Point) Point {
	return Point{p.x + other.x, p.y + other.y}
}
func main() {
	b := new(Bag)
	b.Add(1)
	b.Add(2)

	fmt.Println(b.items)

	p := Property{Value: 1}
	fmt.Println(p.GetValue())
	p.SetValue(2)
	fmt.Println(p.GetValue())

	p1 := Point{1, 2}
	p2 := Point{2, 3}
	p3 := p1.Add(p2)
	fmt.Println(p3)
}
