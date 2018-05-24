package main

import (
	"fmt"
	"time"
)

func main() {
	chanString := make(chan string)
	chanInt := make(chan int)
	chanDone := make(chan bool)

	go func() {
		for {
			select {
			case str := <-chanString:
				fmt.Printf("received a string: %q\n", str)
			case i := <-chanInt:
				fmt.Printf("received an integer: %v\n", i)
			case done := <-chanDone:
				fmt.Printf("received a done signal: %v\n", done)
				return
			}
		}
	}()

	go func() {
		for {
			chanInt <- 1
		}
	}()

	go func() {
		for {
			chanString <- "a"
		}
	}()

	time.Sleep(time.Microsecond * 2)
	chanDone <- true
}
