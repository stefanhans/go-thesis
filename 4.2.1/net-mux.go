package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Create handler to multiplex according to the URL path
	handler := http.NewServeMux()

	// Register anonymous handler function to URL path "/hello"
	handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "I say hello to you!")
	})

	// Register anonymous handler function to URL path "/good-morning"
	handler.HandleFunc("/good-morning", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "I say good morning to you!")
	})

	// Listen and serve on localhost:22365
	log.Fatal(http.ListenAndServe(":22365", handler))
}
