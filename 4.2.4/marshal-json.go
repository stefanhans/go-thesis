package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Animal struct {
	Name         string
	Kind         string
	NumberOfLegs int
}

func main() {
	animals := []Animal{
		Animal{
			"alice",
			"cat",
			4,
		},
		Animal{
			"bob",
			"bird",
			2,
		},
		Animal{
			"curt",
			"fish",
			0,
		},
	}

	// Marshal array of struct
	if byteArray, err := json.Marshal(animals); err != nil {
		log.Fatal(err)
	} else {
		// Unformated JSON
		fmt.Printf("%s\n", byteArray)

		//Formated JSON
		var out bytes.Buffer
		json.Indent(&out, byteArray, "", "\t")
		out.WriteTo(os.Stdout)
	}
}
