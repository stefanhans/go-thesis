package main

import "fmt"

type integer int

func (i integer) add(j integer) integer {
	fmt.Printf("%v + %v = %v\n", i, j, i+j)
	return i + j
}

func main() {
	fmt.Println(integer(1).add(integer(2)).add(integer(3)).add(integer(10)))
}