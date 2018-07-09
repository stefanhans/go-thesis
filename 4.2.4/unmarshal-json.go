package main

import (
	"encoding/json"
	"fmt"
)

type Animal struct {
	Name         string
	Kind         string
	NumberOfLegs int
}

func main() {
	// Define JSON string
	jsonAnimals := []byte(`
	[{"Name":"alice","Kind":"cat","NumberOfLegs":4},
	{"Name":"bob","Kind":"bird","NumberOfLegs":2},
	{"Name":"curt","Kind":"fish","NumberOfLegs":0}]`)

	// Unmarshal JSON into array of struct
	var animals []Animal
	err := json.Unmarshal(jsonAnimals, &animals)
	if err != nil {
		fmt.Println("error:", err)
	}
	// Print type and values
	fmt.Printf("%T: %+v", animals, animals)
}
