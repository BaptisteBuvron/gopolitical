package gopolitical

import (
	"math"
)

type Environment struct {
	Countries             map[string]*Country      `json:"countries"`
	World                 *World                   `json:"world"`
	RelationManager       *RelationManager         `json:"-"`
	Market                *Market                  `json:"market"`
	Agents                map[string]Agent         `json:"-"`
	ConsumptionByHabitant map[ResourceType]float64 `json:"consumptionByHabitant"`
	CurrentDay            int                      `json:"currentDay"`
}

func NewEnvironment(worldWidth int, worldHeight int, countries map[string]*Country, territories []*Territory, prices Prices, consumptionsByHabitant map[ResourceType]float64) *Environment {
	//Map des perceptions que reçoivent les pays à chaque tour
	env := &Environment{
		Countries:             countries,
		World:                 NewWorld(territories, worldWidth, worldHeight),
		RelationManager:       NewRelationManager(),
		Market:                NewMarket(prices),
		ConsumptionByHabitant: consumptionsByHabitant,
		CurrentDay:            1,
	}
	env.Market.Env = env

	// We have only countries as agents
	env.Agents = make(map[string]Agent, len(env.Countries))
	for _, country := range env.Countries {
		env.Agents[country.ID] = country
	}
	return env
}

func (e *Environment) HandleActions(actions []Action) {
	// Permet de ne pas priorisé un agents
	random.Shuffle(len(actions), func(i, j int) {
		actions[i], actions[j] = actions[j], actions[i]
	})
	// Les requêtes s’exécute dans le contexte de l'environnement
	for _, action := range actions {
		action.Execute(e)
	}
}

func (e *Environment) Update() {
	// On fait correspondre les ordres d'achats et de ventes
	e.Market.HandleRequests()

	// Mettre à jour les stocks des territoires à partir des variations
	e.UpdateStocksFromVariation()

	// Mettre à jour les stocks des territoires à partir des consommations des habitants
	e.UpdateStocksFromConsumption()
	e.ApplyRulesOfLife()

	//Add history
	e.UpdateStockHistory()
	e.UpdateMoneyHistory()
	e.UpdateHabitantsHistory()

	e.CurrentDay += 1

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

func (e *Environment) UpdateStockHistory() {
	for _, territory := range e.World.Territories() {
		// copy stock
		copyStock := make(map[ResourceType]float64)
		for k, v := range territory.Stock {
			copyStock[k] = v
		}
		territory.StockHistory[e.CurrentDay] = copyStock
	}
}

func (e *Environment) UpdateMoneyHistory() {
	for _, country := range e.Countries {
		country.MoneyHistory[e.CurrentDay] = country.Money
	}
}

func (e *Environment) UpdateHabitantsHistory() {
	for _, territory := range e.World.Territories() {
		territory.HabitantsHistory[e.CurrentDay] = territory.Habitants
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
