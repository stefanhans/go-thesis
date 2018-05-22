package main

import "fmt"

type Human struct { Name string }

type Man struct { Human }

type Women struct { Human }

type Talker interface { Talk() }

func SpeakOut(t Talker) {
	t.Talk()
}

func (h *Human) Talk() {
	fmt.Printf("I am an human named %s\n", h.Name)
}

func (m *Man) Talk() {
	fmt.Printf("I am a man named %s\n", m.Name)
}

func main() {
	jamie := Human{Name: "Jamie"}
	claude := Man{Human: Human{Name: "Claude"}}
	cameron := Women{Human: Human{Name: "Cameron"}}

	SpeakOut(&jamie)
	SpeakOut(&claude)
	SpeakOut(&cameron)
}
