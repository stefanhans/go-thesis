package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ghodss/yaml"
)

func main() {

	// Read YAML file into slice of bytes
	sliceOfBytes, err := ioutil.ReadFile("family.yaml")
	if err != nil {
		log.Fatal("could not read family.yaml: %v", err)
	}

	// Convert YAML to JSON
	jsonBytes, err := yaml.YAMLToJSON(sliceOfBytes)
	if err != nil {
		fmt.Printf("could not convert YAML to JSON: %v\n", err)
		return
	}
	fmt.Printf("%s\n", jsonBytes)

	// Convert JSON to YAML
	yamlBytes, err := yaml.JSONToYAML(jsonBytes)
	if err != nil {
		fmt.Printf("could not convert JSON to YAML: %v\n", err)
		return
	}
	fmt.Printf("\n%s", yamlBytes)
}
