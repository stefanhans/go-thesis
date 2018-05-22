package main

import "fmt"

type Artist struct { Name string }

type Writer struct { Artist }

type Painter struct { Artist }

func (a *Artist) Talk() {
	fmt.Printf("I am an artist named %s\n", a.Name)
}

func (w *Writer) Talk() {
	fmt.Printf("I am a writer named %s\n", w.Name)
}

func main() {
	artist := Artist{Name: "Alex"}
	writer := Writer{Artist: Artist{Name: "William"}}
	painter := Painter{Artist: Artist{Name: "Paul"}}

	artist.Talk()
	writer.Talk()
	painter.Talk()
}
