package gopolitical

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type Environment struct {
	Countries       map[string]*Country  `json:"-"`
	World           *World               `json:""`
	RelationManager *RelationManager     `json:"-"`
	Market          *Market              `json:"market"`
	RandomGenerator *rand.Rand           `json:"-"`
	lock            sync.Mutex           `json:"-"`
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

	randomSource := rand.NewSource(time.Now().UnixNano())
	env.RandomGenerator = rand.New(randomSource)
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
				responsePercept.Prices = e.Market.Prices
				Respond(fromCountry.In, responsePercept)
				break
			case AttackRequest:
				// TODO: verifier que les pays sont bien voisins au moments de l'attaques
				// On récupère les stocks d'armes des deux pays
				defensiveArmament := req.to.Country.GetTotalStockOf(ARMAMENT)
				offensiveArmament := req.from.GetTotalStockOf(ARMAMENT)
				// On verify que le pays à bien de quoi attaquer
				if offensiveArmament < req.armement {
					log.Printf("Attaque avortée de %v sur %v\n", req.from.ID, req.to.Country.ID)
					continue
				}
				// Il n'utilisera que ce qu'il souhaite
				offensiveArmament = req.armement
				// Si l'attaqué a assez de ressource pour se défendre
				if req.armement < defensiveArmament {
					defensiveArmament = req.armement // Il n'en utilise qu'une partiel
				}

				// Il font une bataille donc il consomme de l'armement
				req.from.Consume(ARMAMENT, offensiveArmament)
				req.to.Country.Consume(ARMAMENT, defensiveArmament)

				// Le taux de réussite correspond à offensif / défensif
				chanceOfFailure := offensiveArmament / defensiveArmament
				log.Printf("%v Attaque %v avec %.2f de réussite\n", req.from.ID, req.to.Country.ID, chanceOfFailure*100)

				// On récupère l'état de la relation actuelle
				relation := e.RelationManager.GetRelation(req.from.ID, req.to.Country.ID)
				attackedCountry := req.to.Country
				if chanceOfFailure < e.RandomGenerator.Float64() { // L'attaque a réussi
					relation = relation / 3
					req.to.TransfertProperty(req.from)
					log.Printf("Attaque réussit\n")
				} else { // L'attaque a échoué
					relation = relation * 2 / 3
					log.Printf("Attaque échoué\n")
				}
				e.RelationManager.UpdateRelation(req.from.ID, attackedCountry.ID, relation)
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
