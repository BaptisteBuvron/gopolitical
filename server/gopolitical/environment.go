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
	for _, country := range e.Countries {
		select {
		case req := <-country.In:
			//Try downcasting to a MarketRessourceRequest
			e.lock.Lock()
			if marketReq, ok := req.(MarketBuyRequest); ok {
				e.Market.handleRequest(marketReq)
			} else {
				log.Println("Une requete n'a pas pu etre traitee")
			}
			e.lock.Unlock()
		default:
		}
	}
}

func (e *Environment) handleMarketRequest(req MarketBuyRequest) {
	e.Market.handleRequest(req)
}
