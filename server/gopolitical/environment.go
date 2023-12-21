package gopolitical

import (
	"fmt"
	"log"
	"sync"
)

type Environment struct {
	Countries   map[string]*Country
	Territories []*Territory
	Market      Market
	wg          *sync.WaitGroup
	lock        sync.Mutex
	Percept     map[string][]Request
}

func NewEnvironment(countries map[string]*Country, territories []*Territory, prices Prices, wg *sync.WaitGroup) Environment {
	//Map des perceptions que recoivent les pays à chaque tour
	percept := make(map[string][]Request)
	for _, country := range countries {
		percept[country.Name] = []Request{}
	}
	return Environment{
		Countries:   countries,
		Territories: territories,
		Market:      NewMarket(prices, percept),
		wg:          wg,
		lock:        sync.Mutex{},
		Percept:     percept,
	}
}

func (e *Environment) Start() {
	log.Println("Start of the environment")
	for {
		e.handleRequests()
	}
}

func (e *Environment) handleRequests() {
	for _, country := range e.Countries {
		select {
		case req := <-country.Out:
			//Try downcasting
			e.lock.Lock()
			switch req := req.(type) {
			case MarketBuyRequest, MarketSellRequest:
				e.Market.handleRequest(req)
				break
			case PerceptRequest:
				fmt.Println("Une requete a ete traitee")
				fromCountry := req.from
				responsePercept := PerceptResponse{events: e.Percept[fromCountry.Name]}
				e.Percept[fromCountry.Name] = []Request{}
				Respond(fromCountry.In, responsePercept)
				break

			default:
				log.Println("Une requete n'a pas pu etre traitee")
			}
			e.lock.Unlock()
		default:
		}
	}
}

func Respond(toChannel Channel, res Request) {
	toChannel <- res
}
