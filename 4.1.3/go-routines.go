package main

import (
	"fmt"
	"time"
)

func count(str string) {
	for i := 0; i < 5; i++ {
		fmt.Printf("%s: %v\n", str, i)
		time.Sleep(time.Microsecond)
	}
}

func main() {
	go count("go count")
	count("count")
}
