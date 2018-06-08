package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func main() {

	// Read YAML file into slice of bytes
	sliceOfBytes, err := ioutil.ReadFile("family.yaml")
	if err != nil {
		log.Fatal("could not read storage: %v", err)
	}

	// Unmarshall unknown YAML
	mapOfUnknown := make(map[interface{}]interface{})

	err = yaml.Unmarshal(sliceOfBytes, &mapOfUnknown)
	if err != nil {
		log.Fatalf("could not unmarshall YAML: %v", err)
	}
	fmt.Printf("%+v\n", mapOfUnknown)
}
