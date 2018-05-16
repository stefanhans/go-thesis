package main

import "fmt"

func main() {
	sliceOfDigits := []string{
		"zero", "one", "two", "three",
		"four", "five", "six", "seven",
		"eight", "nine", "ten"}

	fmt.Printf("sliceOfDigits[0]  :\t%v\n", sliceOfDigits[0])
	fmt.Printf("sliceOfDigits[1:6]:\t%v\n", sliceOfDigits[1:6])
	fmt.Printf("sliceOfDigits[6: ]:\t%v\n", sliceOfDigits[6:])
	fmt.Printf("sliceOfDigits[ :4]:\t%v\n", sliceOfDigits[:4])

}
