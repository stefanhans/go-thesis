package main

import (
	"os"
	"text/template"
)

type animal struct {
	Name         string
	Kind         string
	NumberOfLegs int
}

func main() {

	animals := []animal{
		animal{
			"alice",
			"cat",
			4,
		},
		animal{
			"bob",
			"bird",
			2,
		},
		animal{
			"curt",
			"fish",
			0,
		},
	}

	template, err := template.New("animal").
		Parse("{{.Name}} is a {{.Kind}} with {{.NumberOfLegs}} legs.\n")
	if err != nil {
		panic(err)
	}
	for _, animal := range animals {
		err = template.Execute(os.Stdout, animal)
		if err != nil {
			panic(err)
		}
	}

}
