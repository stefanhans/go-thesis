package main

import "fmt"

type rectangle struct {
	width  int
	height int
}

func main() {
	oneRectangle := rectangle{
		width:  8,
		height: 12,
	}

	fmt.Printf("oneRectangle: %+v\n", oneRectangle)
	fmt.Printf("oneRectangle.width: %v, "+
		"oneRectangle.height: %+v\n",
		oneRectangle.width, oneRectangle.height)
}
