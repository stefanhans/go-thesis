package main

import (
	"fmt"
	"math"
)

type vector []float64

func (v vector) abs() float64 {
	sumOfSquares := 0.0
	for i := range v {
		sumOfSquares += v[i] * v[i]
	}
	return math.Sqrt(sumOfSquares)
}

func main() {
	v := vector{2, 4, 4}
	fmt.Printf("Absolute value of vector %v: %v\n", v, v.abs())
}
