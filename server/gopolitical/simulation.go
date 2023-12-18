package gopolitical

import "sync"

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

	s.wg.Add(len(s.Countries))
	for _, country := range s.Countries {
		go country.Start()
	}
	for {
		//Wait for all agents to finish their actions
		s.wg.Wait()
		//Update the environment
		s.Environment.Update()

		//Restart all agents
		s.wg.Add(len(s.Countries))
	}
}

func (e *Environment) Update() {
	// Add your implementation here
}

func (s *Simulation) Run() {
	go s.Start()
}
