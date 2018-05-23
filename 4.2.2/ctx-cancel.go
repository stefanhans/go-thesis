package main

import (
	"context"
	"fmt"
	"time"
)

func main() {

	// create context and function for cancellation
	ctx, cancel := context.WithCancel(context.Background())

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

	// cancel context and go routines, respectively, and wait for printed messages
	cancel()
	time.Sleep(time.Millisecond)
}
