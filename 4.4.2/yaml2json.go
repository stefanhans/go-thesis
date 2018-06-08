package main

import (
	"fmt"
	"encoding/json"

	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"bytes"
	"os"
)


type Member struct {
	Kind      string `json:"kind"`
	Firstname string `json:"firstname"`
	Age       int    `json:"age"`
}

type Family struct {
	Name    string   `json:"name"`
	Members []Member `json:"members"`
}

func main() {

	// Read YAML file into slice of bytes
	sliceOfBytes, err := ioutil.ReadFile("family.yaml")
	if err != nil {
		log.Fatal("could not read family.yaml: %v", err)
	}


	json, err := yaml.YAMLToJSON(sliceOfBytes)
	if err != nil {
		fmt.Printf("could not convert YAML to JSON: %v\n", err)
		return
	}
	fmt.Printf("%s\n", string(json))

	var f Family
	yaml.Unmarshal(sliceOfBytes, &f)

	fmt.Printf("%s\n", f)


	//Formated JSON
	var out bytes.Buffer
	json.Indent(&out, sliceOfBytes, "", "\t")
	out.WriteTo(os.Stdout)


	yaml, err := yaml.JSONToYAML(json)
	if err != nil {
		fmt.Printf("could not convert JSON to YAML: %v\n", err)
		return
	}
	fmt.Println(string(yaml))
}