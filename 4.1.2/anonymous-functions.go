package main

import (
	"fmt"
)

func main() {

	for i := 0; i < 3; i++ {

		// evaluated anonymous function as closure;
		// variable i is in the outer scope but can be changed from inside the function
		func() {
			fmt.Printf("Variable i: %v\n", i)
			i++
		}()
	}

	fmt.Println()

	for i := 0; i < 3; i++ {

		// evaluated anonymous function with parameter passed by value;
		// variable i is passed by value and only a copy is inside
		func(i int) {
			fmt.Printf("Variable i: %v\n", i)
			i++
		}(i)
	}
}
