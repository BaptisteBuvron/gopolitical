package gopolitical

import (
	"log"
	"sync"
)

type Environment struct {
	Countries             map[string]*Country `json:"-"`
	Territories           []*Territory        `json:"-"`
	Market                *Market             `json:"market"`
	wg                    *sync.WaitGroup
	lock                  sync.Mutex
	Percept               map[string][]Request     `json:"-"`
	ConsumptionByHabitant map[ResourceType]float64 `json:"consumptionByHabitant"`
}

func NewEnvironment(countries map[string]*Country, territories []*Territory, prices Prices, wg *sync.WaitGroup, consumptionsByHabitant map[ResourceType]float64) Environment {
	//Map des perceptions que recoivent les pays à chaque tour
	percept := make(map[string][]Request)
	for _, country := range countries {
		percept[country.Name] = []Request{}
	}
	return Environment{
		Countries:             countries,
		Territories:           territories,
		Market:                NewMarket(prices, percept),
		wg:                    wg,
		lock:                  sync.Mutex{},
		Percept:               percept,
		ConsumptionByHabitant: consumptionsByHabitant,
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
				Respond(country.In, req)
				break
			case PerceptRequest:
				fromCountry := req.from
				responsePercept := PerceptResponse{events: e.Percept[fromCountry.Name]}
				e.Percept[fromCountry.Name] = []Request{}
				Respond(fromCountry.In, responsePercept)
				break

			default:
				log.Println("Une requete n'a pas pu etre traitee")
			}
			//respond to indicate the request was handled
			e.lock.Unlock()
		default:
		}
	}
}

func Respond(toChannel Channel, res Request) {
	toChannel <- res
}

func (e *Environment) UpdateStocksFromVariation() {
	//Mettre à jour les stocks des territoires à partir des variations
	for _, territory := range e.Territories {
		for _, variation := range territory.Variations {
			territory.Stock[variation.Ressource] += variation.Amount
		}
	}
}

func (e *Environment) UpdateStocksFromConsumption() {
	//Mettre à jour les stocks des territoires à partir des consommations
	for _, country := range e.Countries {
		for _, territory := range country.Territories {
			for resource, consumption := range e.ConsumptionByHabitant {
				territory.Stock[resource] -= float64(territory.Habitants) * consumption
			}
		}
	}
}

func (e *Environment) UpdateStockHistory(currentDay int) {
	for _, territory := range e.Territories {
		//copy stock
		copyStock := make(map[ResourceType]float64)
		for k, v := range territory.Stock {
			copyStock[k] = v
		}
		territory.StockHistory[currentDay] = copyStock
	}
}

func (e *Environment) UpdateMoneyHistory(currentDay int) {
	for _, country := range e.Countries {
		country.MoneyHistory[currentDay] = country.Money
	}
}

func (e *Environment) UpdateHabitantsHistory(day int) {
	for _, territory := range e.Territories {
		territory.HabitantsHistory[day] = territory.Habitants
	}
}
