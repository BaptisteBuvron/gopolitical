package gopolitical

import (
	"log"
	"math"
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

func (e *Environment) KillHungryHabitants() {
	totalKilledHabitants := make(map[string]int)
	for _, territory := range e.Territories {
		habitantsHungryByResource := make(map[ResourceType]int)
		for resource, consumption := range e.ConsumptionByHabitant {
			if territory.Stock[resource] < 0 {
				habitantsHungry := math.Ceil(math.Abs(territory.Stock[resource]) / consumption)
				habitantsHungryByResource[resource] = int(habitantsHungry)
			}
		}
		//get max habitants hungry
		maxHabitantsHungry := 0
		for _, habitantsHungry := range habitantsHungryByResource {
			if habitantsHungry > maxHabitantsHungry {
				maxHabitantsHungry = habitantsHungry
			}
		}
		//On tue la moitié des habitants qui ont faim
		killedHabitants := int(math.Ceil(float64(maxHabitantsHungry) / 2))
		if territory.Habitants-killedHabitants <= 0 {
			killedHabitants = territory.Habitants - 1
		}
		territory.Habitants -= killedHabitants
		totalKilledHabitants[territory.Country.Name] += killedHabitants
	}
	for countryName, killedHabitants := range totalKilledHabitants {
		log.Println("Famine : ", killedHabitants, " habitants sont mort de faim ", countryName)
	}

}

func (e *Environment) BirthHabitants() {
	totalBirthHabitants := make(map[string]int)
	for _, territory := range e.Territories {
		birth := int(math.Ceil(float64(territory.Habitants) * 0.02))
		territory.Habitants += birth
		totalBirthHabitants[territory.Country.Name] += birth
	}
	for countryName, birthHabitants := range totalBirthHabitants {
		log.Println("Naissance : ", birthHabitants, " habitants de ", countryName)
	}
}
