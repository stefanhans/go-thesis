package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// create context with timeout for cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// call go routines waiting for its cancellation
	for i := 0; i < 4; i++ {
		go func(i int) {
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("%v: ctx.Done(): %v\n", i, ctx.Err())
					return
				}
			}
		}(i)
	}

	// wait for deadline cancelling context and go routines, respectively
	// and for printed messages
	time.Sleep(time.Second)
}
