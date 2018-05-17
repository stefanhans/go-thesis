package main

import "fmt"

type Author struct {
	Name string
	City string
}

type Book struct {
	Title string
	Author Author
}

func main() {

	b := Book{
		Title: "Just A Book",
		Author: Author{
			Name: "W. Riter",
			City: "New York",
		},
	}

	fmt.Printf("%q was written from %s at %s\n",
		b.Title, b.Author.Name, b.Author.City)
}
