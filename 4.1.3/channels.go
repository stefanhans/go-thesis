package main

import "fmt"

func main() {

	channel := make(chan int)

	go func() { channel <- 41 }()

	fmt.Printf("Received: %v\n", <-channel)
}
