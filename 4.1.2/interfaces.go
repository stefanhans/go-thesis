package main

import (
	"math"
	"fmt"
)

type rectangle struct { width, height float64 }

type circle struct { radius float64 }

type calculator interface {
	calculateArea() float64
}

func (r rectangle) calculateArea() float64 {
	return r.width * r.height
}

func (c circle) calculateArea() float64 {
	return math.Pi * c.radius * c.radius
}

func showArea(c calculator) {
	fmt.Printf("%+v: area: %v\n",
		c, c.calculateArea())
}

func main() {
	r := rectangle{width: 3, height: 4}
	c := circle{radius: 5}
	showArea(r)
	showArea(c)
}