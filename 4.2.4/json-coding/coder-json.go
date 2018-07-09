package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

func main() {

	// Define decoder for reading JSON string
	decoder := json.NewDecoder(os.Stdin)

	// Define encoder for outputting JSON
	encoder := json.NewEncoder(os.Stdout)

	for {
		// Decode string into map
		var jsonMap map[string]interface{}
		if err := decoder.Decode(&jsonMap); err != nil {
			// EOF expected
			return
		}
		// Range map to capitalize string values
		for key := range jsonMap {
			if convertedValue, ok := jsonMap[key].(string); ok {
				jsonMap[key] = strings.Title(convertedValue)
			}
		}

		// Encode output
		if err := encoder.Encode(&jsonMap); err != nil {
			log.Println(err)
		}
	}
}
