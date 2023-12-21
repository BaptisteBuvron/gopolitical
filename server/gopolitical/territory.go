package gopolitical

type Territory struct {
	X          int
	Y          int
	Variations []Variation
	Stock      map[ResourceType]int
	Country    *Country
}

func NewTerritory(x int, y int, variations []Variation, stock map[ResourceType]int, country *Country) *Territory {
	return &Territory{x, y, variations, stock, country}
}

func (c Territory) Start() {

}

func (c Territory) Percept() {

}

func (c Territory) Deliberate() {

}

func (c Territory) Act() {

}
