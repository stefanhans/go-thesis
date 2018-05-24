package main

import (
	"fmt"
	"math"
)

type rectangle struct{ width, height float64 }

type circle struct{ radius float64 }

// interface definition
type calculator interface {
	calculateArea() float64
}

// implement interface as method
func (r rectangle) calculateArea() float64 {
	return r.width * r.height
}

// implement interface as method
func (c circle) calculateArea() float64 {
	return math.Pi * c.radius * c.radius
}

// interface as function parameter
func showArea(c calculator) {
	fmt.Printf("%+v: area: %v\n",
		c, c.calculateArea())
}

func main() {
	r := rectangle{width: 3, height: 4}
	c := circle{radius: 5}

	// interface implementations as function parameter
	showArea(r)
	showArea(c)

	// interface for variable declaration
	var calc calculator
	calc = r
	fmt.Printf("%+v: area: %v\n", calc, calc.calculateArea())
	calc = c
	fmt.Printf("%+v: area: %v\n", calc, calc.calculateArea())
}
