package main

import (
	"encoding/json"
	"fmt"
)

type Product struct {
	Name     string
	Price    float64
	Quantity int
}

func toJson(p Product) (string, error) {
	jsonData, error := json.Marshal(p)
	if error != nil {
		return "", error
	}
	return string(jsonData), error
}

func fromJson(jsonStr string) (Product, error) {
	var p Product
	error := json.Unmarshal([]byte(jsonStr), &p)
	if error != nil {
		return Product{}, error
	}
	return p, nil
}
func main() {
	candy := Product{
		Name:     "haribo",
		Price:    50.15,
		Quantity: 3,
	}
	jsonStr, error := toJson(candy)
	if error != nil {
		fmt.Println("FAIL:", error)
	} else {
		fmt.Println("JSON:", jsonStr)
	}
	jsonTest := `{"Name":"Alibaba","Price":19.99,"Quantity":6}`
	testProd, error := fromJson(jsonTest)

	if error != nil {
		fmt.Println("FAIL", error)
	} else {
		fmt.Println("Product is ", testProd)
	}
	fmt.Print()
}
