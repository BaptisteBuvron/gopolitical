package gopolitical

import (
	"fmt"
	"sync"
	"time"
)

type Simulation struct {
	SecondByDay float64
	Environment Environment
	Territories []Territory
	Countries   map[string]Country
	wg          *sync.WaitGroup
}

func NewSimulation(
	secondByDay float64,
	prices Prices,
	countries map[string]Country,
	territories []Territory,
	wg *sync.WaitGroup,
) Simulation {
	return Simulation{secondByDay, NewEnvironment(countries, territories, prices, wg), territories, countries, wg}
}

func (s *Simulation) Start() {
	//Launch all agents and added a channel to the environment

	fmt.Println("Start of the simulation : ")
	fmt.Println("Number of countries : ", len(s.Countries))

	for {
		//Restart all agents
		fmt.Println("Start of a new day")
		for _, country := range s.Countries {
			s.wg.Add(1)
			go country.Start()
		}
		//Wait for all agents to finish their actions
		s.wg.Wait()
		fmt.Println("End of the day")
		//Update the environment
		s.Environment.Update()

		//Wait the other day
		time.Sleep(time.Duration(s.SecondByDay) * time.Second)
	}
}

func (e *Environment) Update() {
	// Add your implementation here
}

func (s *Simulation) Run() {
	go s.Start()
}
