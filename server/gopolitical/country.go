package gopolitical

import (
	"fmt"
	"sync"
)

type Country struct {
	Agent
	Color       string
	Territories []Territory
	Money       float64
	wg          *sync.WaitGroup
	In          Channel
	Out         Channel
}

func NewCountry(id string, name string, color string, territories []Territory, money float64, wg *sync.WaitGroup, in Channel, out Channel) Country {
	return Country{Agent{id, name}, color, territories, money, wg, in, out}
}

func (c Country) Start() {
	fmt.Printf("Country %s started\n", c.Name)
	c.Percept()
	requests := c.Deliberate()
	c.Act(requests)
	c.wg.Done()
}

func (c Country) Percept() {
	perceptRequest := PerceptRequest{from: c}
	c.Out <- perceptRequest
	perceptResponse := <-c.In

	//Downcast to a PerceptResponse
	if perceptResponse, ok := perceptResponse.(PerceptResponse); ok {
		fmt.Println(len(perceptResponse.events))
	}

	fmt.Printf("Country %s percept\n", c.Name)
}

func (c Country) Deliberate() []Request {
	fmt.Printf("Country %s deliberate\n", c.Name)

	//Si le pays a plus de ressources que ce qu'il lui faut, il vend le surplus

	return nil
}

func (c Country) Act(requests []Request) {
	fmt.Printf("Country %s act\n", c.Name)
}
