package gopolitical

type Country struct {
	Agent
	Color       string
	Territories []Territory
	Money       float64
}

func NewCountry(id string, name string, color string, territories []Territory, money float64) Country {
	return Country{Agent{id, name}, color, territories, money}
}

func (c Country) Start() {

}

func (c Country) Percept() {

}

func (c Country) Deliberate() {

}

func (c Country) Act() {

}
