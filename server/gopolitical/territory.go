package gopolitical

import "fmt"

type Territory struct {
	Agent
	X          int
	Y          int
	Variations []Variation
	Country    Country
}

func NewTerritory(x int, y int, variations []Variation, country Country) Territory {
	id := fmt.Sprintf("%d-%d", x, y)
	return Territory{Agent{id, id}, x, y, variations, country}
}

func (c Territory) Start() {

}

func (c Territory) Percept() {

}

func (c Territory) Deliberate() {

}

func (c Territory) Act() {

}
