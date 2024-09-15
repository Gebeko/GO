package main

import "fmt"

type Shape interface {
	Area() float64
}

type Circle struct {
	Radius float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (c Circle) Area() float64 {
	area := c.Radius * c.Radius * 3.142
	return area
}

func (r Rectangle) Area() float64 {
	area := r.Height * r.Width
	return area
}

func PrintArea(s Shape) {
	fmt.Printf("Area = %v\n", s.Area())
}

func main() {
	circle := Circle{
		Radius: 3,
	}
	rectangle := Rectangle{
		Width:  3.3,
		Height: 5.5,
	}
	PrintArea(circle)
	PrintArea(rectangle)
}
