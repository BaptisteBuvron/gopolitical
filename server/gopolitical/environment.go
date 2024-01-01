package gopolitical

import (
	"log"
	"sync"
)

type Environment struct {
	Countries       map[string]*Country `json:"-"`
	World           *World              `json:""`
	RelationManager *RelationManager    `json:"-"`
	Market          *Market             `json:"market"`
	lock            sync.Mutex
	Percept         map[string][]Request `json:"-"`
}

func NewEnvironment(worldWidth int, worldHeight int, countries map[string]*Country, territories []*Territory, prices Prices) *Environment {
	// Map des perceptions que reçoivent les pays à chaque tour
	percept := make(map[string][]Request)
	for _, country := range countries {
		percept[country.Name] = []Request{}
	}
	env := &Environment{
		Countries:       countries,
		World:           NewWorld(territories, worldWidth, worldHeight),
		RelationManager: NewRelationManager(),
		Market:          NewMarket(prices, percept),
		lock:            sync.Mutex{},
		Percept:         percept,
	}
	env.Market.Env = env
	return env
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
				fromCountry := req.from
				responsePercept := PerceptResponse{events: e.Percept[fromCountry.Name]}
				e.Percept[fromCountry.Name] = []Request{}
				responsePercept.RelationManager = e.RelationManager
				responsePercept.World = e.World
				Respond(fromCountry.In, responsePercept)
				break
			case AttackRequest:
				log.Printf("Attaque %v\n", req)
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

func (e *Environment) UpdateStocksFromVariation() {
	// Mettre à jour les stocks des territoires à partir des variations
	for _, territory := range e.World.Territories() {
		for _, variation := range territory.Variations {
			territory.Stock[variation.Ressource] += variation.Amount
		}
	}
}

func (e *Environment) UpdateStocksFromConsumption() {
	// Mettre à jour les stocks des territoires à partir des consommations
	for _, country := range e.Countries {
		for _, territory := range country.Territories {
			foodConsumption := float64(territory.Habitants) * FOOD_BY_HABITANT
			territory.Stock["food"] -= foodConsumption

			waterConsumption := float64(territory.Habitants) * WATER_BY_HABITANT
			territory.Stock["water"] -= waterConsumption
		}
	}
}

func (e *Environment) UpdateStockHistory(currentDay int) {
	for _, territory := range e.World.Territories() {
		// copy stock
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
	for _, territory := range e.World.Territories() {
		territory.HabitantsHistory[day] = territory.Habitants
	}
}
