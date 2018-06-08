package main

import (
	"fmt"
	"io/ioutil"
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

	// Read YAML file into slice of bytes
	sliceOfBytes, err := ioutil.ReadFile("family.yaml")
	if err != nil {
		log.Fatal("could not read storage: %v", err)
	}

	// Unmarshall slice of bytes
	var family Family
	err = yaml.Unmarshal(sliceOfBytes, &family)
	if err != nil {
		panic(err)
	}

	// Print Family struct
	fmt.Printf("%+v\n", family)
}
