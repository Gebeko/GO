package main

import "fmt"

func main() {
	var a int = 32
	var b float64
	b = 1.14
	var c string = "Hello"
	var d bool = false

	e := 15
	f := "World"

	fmt.Printf("Integer a is %v and integer e is %v. Both are a %T\n", a, e, a)
	fmt.Printf("Float a is %v and is a %T\n", b, b)
	fmt.Printf("Did i close the door? %v \n", d)
	fmt.Printf("%v,%v", c, f)
}
