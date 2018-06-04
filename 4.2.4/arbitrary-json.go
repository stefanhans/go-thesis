package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	// Define JSON string
	b := []byte(`{"Kind":"Mouse","NumberOfLegs":4,"Names":["Bernard","Bianca"]}`)

	// Unmarshal JSON into map of empty interfaces
	var mapOfInterfaces map[string]interface{}
	err := json.Unmarshal(b, &mapOfInterfaces)
	if err != nil {
		fmt.Println("error:", err)
	}

	// Range with type switch
	for key, value := range mapOfInterfaces {
		switch convertedValue := value.(type) {
		case string:
			fmt.Printf("%q (string): %s\n", key, convertedValue)
		case float64:
			fmt.Printf("%q (float64): %f\n", key, convertedValue)
		case []interface{}:
			fmt.Printf("%q ([]interface{}): %v\n", key, convertedValue)
		}
	}
}
