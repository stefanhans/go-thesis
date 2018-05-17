package main

import "fmt"

type rectangle struct {
	width  int
	height int
}
func (rect *rectangle) area() int {
	return rect.width * rect.height
}

func main() {
	quadrat := rectangle{
		width:  8,
		height: 8,
	}

	fmt.Printf("quadrat.area(): %+v\n", quadrat.area())
}
