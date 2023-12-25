package gopolitical

import (
	"fmt"
	"sync"
	"time"
)

type Simulation struct {
	SecondByDay float64             `json:"secondByDay"`
	Environment Environment         `json:"environment"`
	Territories []*Territory        `json:"territories"`
	Countries   map[string]*Country `json:"countries"`
	wg          *sync.WaitGroup
	WebSocket   *WebSocket `json:"-"`
}

const (
	WATER_BY_HABITANT = 0.5
	FOOD_BY_HABITANT  = 0.5
)

func NewSimulation(
	secondByDay float64,
	prices Prices,
	countries map[string]*Country,
	territories []*Territory,
	wg *sync.WaitGroup,
) Simulation {
	return Simulation{secondByDay, NewEnvironment(countries, territories, prices, wg), territories, countries, wg, nil}
}

func (s *Simulation) Start() {
	//Launch all agents and added a channel to the environment

	s.WebSocket = NewWebSocket(s)
	go s.WebSocket.Start()

	fmt.Println("Start of the simulation : ")
	fmt.Println("Number of countries : ", len(s.Countries))
	fmt.Println("Number of territories : ", len(s.Territories))

	for _, country := range s.Countries {
		fmt.Println("Nombre de territoires dans  : ", country.Name, " : ", len(country.Territories))
	}

	go s.Environment.Start()

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
		//Mettre à jour les stocks des territoires à partir des variations
		s.Environment.UpdateStocksFromVariation()
		//Mettre à jour les stocks des territoires à partir des consommations des habitants
		s.Environment.UpdateStocksFromConsumption()

		//On fait corrrespondre les ordres d'achats et de ventes
		s.Environment.Market.HandleRequests()

		//Wait the other day
		time.Sleep(time.Duration(s.SecondByDay) * time.Second)
		//Send update to the websocket
		s.WebSocket.SendUpdate()
	}
}
