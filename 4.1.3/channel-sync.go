package main

import (
	"fmt"
)

func main() {
	done := make(chan bool)

	go func() {
		fmt.Println("do some work until it's done")
		done <- true
	}()

	fmt.Println("without channel it would be terminated here")
	<-done
}
