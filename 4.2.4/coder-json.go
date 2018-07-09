package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

func main() {

	// Define JSON string
	const jsonStream = `{"Name":"alice","Kind":"cat","NumberOfLegs":"4"}
					{"Name":"bob","Kind":"bird","NumberOfLegs":"2"}
					{"Name":"curt","Kind":"fish","NumberOfLegs":"0"}`

	// Define decoder for reading JSON string
	decoder := json.NewDecoder(strings.NewReader(jsonStream))

	// Define encoder for outputting JSON
	encoder := json.NewEncoder(os.Stdout)

	// Until EOF
	for {
		// Decode string into map
		var jsonMap map[string]interface{}
		if err := decoder.Decode(&jsonMap); err != nil {
			return
		}
		// Range map to capitalize string values
		for key := range jsonMap {
			switch convertedValue := jsonMap[key].(type) {
			case string:
				jsonMap[key] = strings.Title(convertedValue)
			}
		}

		// Encode output
		if err := encoder.Encode(&jsonMap); err != nil {
			log.Println(err)
		}
	}
}
