package gopolitical

type Territory struct {
	X          int
	Y          int
	Variations []Variation
	Country    Country
}

func NewTerritory(x int, y int, variations []Variation, country Country) Territory {
	return Territory{x, y, variations, country}
}

func (c Territory) Start() {

}

func (c Territory) Percept() {

}

func (c Territory) Deliberate() {

}

func (c Territory) Act() {

}
