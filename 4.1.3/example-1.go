package main

import "fmt"

func main() {
	messages := make(chan string)

	go func() { messages <- "I came from a channel..." }()

	fmt.Println(<-messages)
}
