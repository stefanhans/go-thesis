package main

import (
	"fmt"
	"net/http"
)

// function to handle a request, i.e. write a reply
func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("request %q received from : %v\n", r.URL, r.RemoteAddr)

	fmt.Fprintf(w, "Hello, %s.", r.URL.Path[1:])
}

func main() {

	// listen on TCP for localhost on port 8080
	// and serve the handled request on incoming connections
	http.ListenAndServe("localhost:8080",
		http.HandlerFunc(handler))
}
