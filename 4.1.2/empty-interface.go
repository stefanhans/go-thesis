package main

import (
	"errors"
	"fmt"
)

func processExpectedTypes(in interface{}) (err error) {
	switch value := in.(type) {
	case string:
		fmt.Printf("Processing the string: %q\n", value)
	case int:
		fmt.Printf("Processing the integer: %d\n", value)
	default:
		return errors.New(fmt.Sprintf("no expected type: %T", in))
	}
	return nil
}

func main() {
	for _, value := range []interface{}{1, "2", 3.0} {
		if err := processExpectedTypes(value); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
