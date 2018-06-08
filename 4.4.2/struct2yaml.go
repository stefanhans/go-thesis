package main

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

type Member struct {
	Kind      string `yaml:"kind"`
	Firstname string `yaml:"firstname"`
	Age       int    `yaml:"age"`
}

type Family struct {
	Name    string   `yaml:"name"`
	Members []Member `yaml:"members"`
}

func main() {

	// Define new family struct
	family := Family{
		Name: "Cook",
		Members: []Member{
			Member{
				Kind: "Daddy",
				Firstname: "Sam",
				Age: 31,
			},
			Member{
				Kind: "Mummy",
				Firstname: "Patricia",
				Age: 28,
			},
			Member{
				Kind: "Son",
				Firstname: "Peter",
				Age: 3,
			},
		},
	}

	// Marshall struct into slice of bytes
	sliceOfBytes, err := yaml.Marshal(&family)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Print YAML string
	fmt.Printf("%s", string(sliceOfBytes))
}
