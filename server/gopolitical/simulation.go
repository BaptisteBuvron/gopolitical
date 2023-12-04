package gopolitical

type Simulation struct {
	SecondByDay float64
	Environment Environment
	Territories []Territory
	Countries   map[string]Country
}

func NewSimulation(
	secondByDay float64,
	prices Prices,
	countries map[string]Country,
	territories []Territory,
) Simulation {
	agt := make(map[string]AgentI)
	for _, country := range countries {
		agt[country.ID] = country
	}
	return Simulation{secondByDay, NewEnvironment(agt, prices), territories, countries}
}

func (s *Simulation) Start() {
	//Launch all agents
	//
}
