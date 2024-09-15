package main

import "fmt"

func main() {
	// 1st task
	fmt.Println("Please enter 2 numbers!")
	num0 := 0
	num1 := 0
	fmt.Scan(&num0, &num1)
	sum := num0 + num1
	fmt.Printf("The sum is %v\n", sum)
	// 2nd task
	fmt.Println("Enter 2 strings!")
	str1, str2 := "", ""
	fmt.Scan(&str1, &str2)
	fmt.Printf("Reversed is %v,%v\n", str2, str1)
	// 3rd task
	fmt.Println("Enter 2 numbers for some math!")
	num2, num3 := 0, 0
	fmt.Scan(&num2, &num3)
	quo := num2 / num3
	rem := num2 % num3
	fmt.Printf("The quitotent is %v, and the remainder is %v", quo, rem)
}
