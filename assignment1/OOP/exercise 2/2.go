package main

import "fmt"

type Employee struct {
	Name string
	ID   int
}

type Manager struct {
	Employee
	Department string
}

func (e Employee) Work() {
	fmt.Printf("Employee -> Name: %v, ID: %v\n", e.Name, e.ID)
}
func main() {
	person := Manager{
		Employee: Employee{
			Name: "Gebek",
			ID:   1,
		},
		Department: "Human Resource",
	}
	person.Work()
}
