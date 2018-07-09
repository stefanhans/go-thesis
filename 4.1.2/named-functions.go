package main

import (
	"errors"
	"fmt"
)

// named variadic function with multiple return values
func divide(divisor int, dividends ...int) (float64, error) {
	if divisor == 0 {
		return 0, errors.New("Division by zero")
	} else {
		dividend := 0
		for _, d := range dividends {
			dividend += d
		}
		return float64(dividend / divisor), nil
	}
}

func main() {

	// named variadic function call returning multiple values
	quotient, err := divide(2, 1, 2, 3)

	if err != nil {
		fmt.Printf("Division failed: %v\n", err)
	} else {
		fmt.Printf("Quotient is %v\n", quotient)
	}

	// division by zero
	quotient, err = divide(0, 1, 2, 3)

	if err != nil {
		fmt.Printf("Division failed: %v\n", err)
	} else {
		fmt.Printf("Quotient is %v\n", quotient)
	}

	// empty variable parameter uses type's default value, i.e. 0
	quotient, err = divide(2)

	if err != nil {
		fmt.Printf("Division failed: %v\n", err)
	} else {
		fmt.Printf("Quotient is %v\n", quotient)
	}
}
