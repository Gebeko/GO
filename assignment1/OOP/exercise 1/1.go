package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func (obj Person) Greet() {
	fmt.Printf("Greetings, I am a human being not a robot. My name is %v, and my age is %v\n", obj.Name, obj.Age)
}

func main() {
	Baatar := Person{
		Name: "Gebek",
		Age:  22,
	}
	Baatar.Greet()
}
