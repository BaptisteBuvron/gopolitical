package gopolitical

import (
	"log"
	"sync"
)

type Environment struct {
	Countries   map[string]Country
	Territories []Territory
	Market      Market
	wg          *sync.WaitGroup
	lock        sync.Mutex
}

func NewEnvironment(countries map[string]Country, territories []Territory, prices Prices, wg *sync.WaitGroup) Environment {
	return Environment{
		Countries:   countries,
		Territories: territories,
		Market:      NewMarket(prices),
		wg:          wg,
		lock:        sync.Mutex{},
	}
}

func (e *Environment) Start() {
	log.Println("Start of the environment")
	for {
		e.handleRequests()
	}
}

func (e *Environment) handleMarketRequest(req MarketBuyRequest) {
	e.Market.handleRequest(req)
}

func (e *Environment) handleRequests() {
	for _, country := range e.Countries {
		select {
		case req := <-country.Out:
			//Try downcasting to a MarketRessourceRequest
			e.lock.Lock()
			switch req := req.(type) {
			case MarketBuyRequest:
				e.handleMarketRequest(req)
			case PerceptRequest:
				fromCountry := req.from
				responsePercept := PerceptResponse{events: []Request{}}
				fromCountry.In <- responsePercept
			default:
				log.Println("Une requete n'a pas pu etre traitee")
			}
			e.lock.Unlock()
		default:
		}
	}
}
