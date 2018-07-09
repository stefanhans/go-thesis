package main

import "fmt"

type Artist struct {
	Name string
	Art  string
}

type Human interface {
	Talk()
}

func (p *Artist) Talk() {
	fmt.Printf("I am only an artist named %s\n", p.Name)
}

type Writer struct {
	Country string
	Artist
}

func (p *Writer) Talk() {
	fmt.Printf("I am a writer named %s\n", p.Name)
}

type Painter struct {
	Country string
	Artist
}

func (p *Painter) Talk() {
	fmt.Printf("I am a painter named %s\n", p.Name)
}

func IAm(h Human) {
	h.Talk()
}

func main() {
	artist := Artist{Name: "Tom", Art: "artist"}
	painter := Painter{Artist: Artist{Name: "Marc", Art: "painter"}, Country: "France"}
	writer := Writer{Artist: Artist{Name: "Edgar", Art: "writer"}, Country: "USA"}

	IAm(&artist)
	IAm(&painter)
	IAm(&writer)
}
