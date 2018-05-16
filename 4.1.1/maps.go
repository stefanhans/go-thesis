package main

import "fmt"

func main() {
	mapOfDigits := make(map[int]string)

	for i, word := range []string{
		"zero", "one", "two", "three",
		"four", "five", "six", "seven",
		"eight", "nine", "ten"} {
		mapOfDigits[i] = word
	}

	fmt.Println(mapOfDigits)
}
