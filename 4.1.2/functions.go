package main

import (
	"fmt"
	"errors"
)

func divide(dividend, divisor int) (float64, error) {
	if divisor == 0 {
		return 0, errors.New("Division by zero")
	} else {
		return float64(dividend/divisor), nil
	}
}

// func() is a function type
// func(){} is a function literal which represents an anonymous function
// func(){}() is a used anonymous function or a used evaluated function literal

func main() {
	if quotient, err := divide(10, 0); err != nil {
		fmt.Printf("Division failed: %v\n", err)
	} else {
		fmt.Printf("Quotient is %v\n", quotient)
	}

	var f func()

	f = func() {}
	f()

	for i := 0; i < 5; i++ {
		func() {
			fmt.Println("What happens here? ", i)
			i++
		}()
	}

	//
	for j := 0; j < 5; j++ {
		func(j int) {
			fmt.Println("What happens here? ", j)
			j++
		}(j)
	}
}
