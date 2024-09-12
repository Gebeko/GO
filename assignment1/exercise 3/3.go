package main

import "fmt"

var num int

func main() {
	// 1st task
	fmt.Printf("Enter a number: ")
	fmt.Scan(&num)

	if num > 0 {
		fmt.Printf("This number is positive\n")
	} else if num < 0 {
		fmt.Printf("This number is negative\n")
	} else {
		fmt.Printf("This number is 0\n")

	}
	// 2nd task
	var sum int = 0

	for i := 1; i <= 10; i++ {
		sum += i
	}
	fmt.Printf("The sum of first 10 natural number is %v\n", sum)
	// 3rd task
	var day int
	fmt.Printf("Enter a number (1 to 7 for the day): \n")
	fmt.Scan(&day)

	switch day {
	case 1:
		fmt.Printf("Monday")
	case 2:
		fmt.Printf("Tuesday")
	case 3:
		fmt.Printf("Wednesday")
	case 4:
		fmt.Printf("Thursday")
	case 5:
		fmt.Printf("Friday")
	case 6:
		fmt.Printf("Saturday")
	case 7:
		fmt.Printf("Sunday")
	default:
		fmt.Printf("Invalid input, please enter a number between 1 and 7.")
	}

}
