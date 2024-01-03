package gopolitical

import (
	"math"
	"math/rand"
	"time"
)

type Environment struct {
	Countries             map[string]*Country      `json:"countries"`
	World                 *World                   `json:"world"`
	RelationManager       *RelationManager         `json:"-"`
	Market                *Market                  `json:"market"`
	RandomGenerator       *rand.Rand               `json:"-"`
	Percept               map[string][]Request     `json:"-"`
	ConsumptionByHabitant map[ResourceType]float64 `json:"consumptionByHabitant"`
}

func NewEnvironment(worldWidth int, worldHeight int, countries map[string]*Country, territories []*Territory, prices Prices, consumptionsByHabitant map[ResourceType]float64) *Environment {
	//Map des perceptions que reçoivent les pays à chaque tour
	percept := make(map[string][]Request)
	for _, country := range countries {
		percept[country.Name] = []Request{}
	}
	env := &Environment{
		Countries:             countries,
		World:                 NewWorld(territories, worldWidth, worldHeight),
		RelationManager:       NewRelationManager(),
		Market:                NewMarket(prices, percept),
		Percept:               percept,
		ConsumptionByHabitant: consumptionsByHabitant,
	}
	env.Market.Env = env

	randomSource := rand.NewSource(time.Now().UnixNano())
	env.RandomGenerator = rand.New(randomSource)
	return env
}

func (e *Environment) Start() {
	Debug("Environment", "Start")
	for {
		e.handleRequests()
	}
}

func (e *Environment) handleRequests() {
	// No need to lock
	for _, country := range e.Countries {
		select {
		case req := <-country.Out:
			//Try downcasting
			switch req := req.(type) {
			case MarketBuyRequest, MarketSellRequest:
				e.Market.handleRequest(req)
				Respond(country.In, req)
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
				// TODO: Pas grave si il se trompe, il y aura du tire allié
				// On récupère les stocks d'armes des deux pays
				defensiveArmament := req.to.Country.GetTotalStockOf(ARMAMENT)
				offensiveArmament := req.from.GetTotalStockOf(ARMAMENT)
				// On vérifie que le pays à bien de quoi attaquer
				if offensiveArmament < req.armement && req.armement > 0 && offensiveArmament > 0 {
					Debug(req.from.Name, "[Environment] Attaque avortée sur %s (%s)\n", req.to.Name, req.to.Country.Name)
					Respond(req.from.In, AttackResponse{})
					continue
				}
				// Il n'utilisera que ce qu'il souhaite
				offensiveArmament = req.armement
				// Si l'attaqué a assez de ressource pour se défendre
				if defensiveArmament < 0 {
					defensiveArmament = 0
				}
				if req.armement < defensiveArmament {
					defensiveArmament = req.armement // Il n'en utilise qu'une partie
				}

				// Il font une bataille donc il consomme de l'armement
				req.from.Consume(ARMAMENT, offensiveArmament)
				req.to.Country.Consume(ARMAMENT, defensiveArmament)

				// Le taux de réussite correspond à offensif / défensif avec un bonus de 10%
				chanceOfCapture := 1 - (offensiveArmament/defensiveArmament)*0.8
				Debug(req.from.Name, "Attaque %v (%v) avec %.0f%% de réussite", req.to.Name, req.to.Country.Name, chanceOfCapture*100)

				// On récupère l'état de la relation actuelle
				relation := e.RelationManager.GetRelation(req.from.ID, req.to.Country.ID)
				attackedCountry := req.to.Country
				if chanceOfCapture > e.RandomGenerator.Float64() { // L'attaque a réussi
					relation = relation / 3
					req.to.TransfertProperty(req.from)
					Debug("Environment", "Capturé !")
				} else { // L'attaque a échoué
					relation = relation * 2 / 3
					Debug("Environment", "Échec !")
				}
				e.RelationManager.UpdateRelation(req.from.ID, attackedCountry.ID, relation)
				Respond(req.from.In, AttackResponse{})
				break
			default:
				Debug("Environment", "Une requête n'a pas pu être traitée")
			}
			//respond to indicate the request was handled
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
			for resource, consumption := range e.ConsumptionByHabitant {
				territory.Stock[resource] -= float64(territory.Habitants) * consumption
			}
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

func (e *Environment) ApplyRulesOfLife() {
	totalKilledHabitants := make(map[string]int)
	for _, territory := range e.World.Territories() {
		habitantsHungryByResource := make(map[ResourceType]int)
		for resource, consumption := range e.ConsumptionByHabitant {
			quantity := territory.Stock[resource]
			if quantity < 0 {
				habitantsHungry := math.Ceil(math.Abs(quantity) / consumption)
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
		//On tue un dixième des habitants qui ont faim
		killedHabitants := int(math.Ceil(float64(maxHabitantsHungry) * STARVATION_RATIO))
		if territory.Habitants-killedHabitants <= 0 {
			killedHabitants = territory.Habitants - 1
		}
		territory.Habitants -= killedHabitants
		totalKilledHabitants[territory.Country.Name] += killedHabitants
		if killedHabitants == 0 {
			birth := int(math.Ceil(float64(territory.Habitants) * BIRTH_RATIO))
			territory.Habitants += birth
		}
	}
	for countryName, killedHabitants := range totalKilledHabitants {
		Debug("Environment", "[%s] %d habitants sont mort de faim", countryName, killedHabitants)
	}
}
