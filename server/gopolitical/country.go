package gopolitical

import "sync"

type Country struct {
	Agent
	Color       string
	Territories []Territory
	Money       float64
	wg          *sync.WaitGroup
	In          Channel
	Out         Channel
}

func NewCountry(id string, name string, color string, territories []Territory, money float64, wg *sync.WaitGroup) Country {
	return Country{Agent{id, name}, color, territories, money, wg}
}

func (c Country) Start() {
	for {
		c.Percept()
		c.Deliberate()
		c.Act()
		c.wg.Done()
		c.wg.Wait()
	}
}

func (c Country) Percept() {

}

func (c Country) Deliberate() {

}

func (c Country) Act() {

}
