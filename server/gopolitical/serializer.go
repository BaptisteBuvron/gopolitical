package gopolitical

import (
	"encoding/json"
	"os"
	"sync"
)

type PartialResource struct {
	Name  ResourceType `json:"name"`
	Price float64      `json:"price"`
}

type PartialCountry struct {
	Name  string  `json:"name"`
	ID    string  `json:"id"`
	Color string  `json:"color"`
	Money float64 `json:"money"`
}

type PartialRelation struct {
}

type PartialTerritory struct {
	X          int                `json:"x"`
	Y          int                `json:"y"`
	Country    string             `json:"country"`
	Habitants  int                `json:"habitants"`
	Name       string             `json:"name"`
	Variations []PartialVariation `json:"variations"`
	Stocks     []PartialStock     `json:"stocks"`
}

type PartialStock struct {
	Resource ResourceType `json:"name"`
	Amount   float64      `json:"value"`
}

type PartialVariation struct {
	Name  ResourceType `json:"name"`
	Value float64      `json:"value"`
}

type PartialConsumptionByHabitant struct {
	Resource ResourceType `json:"name"`
	Amount   float64      `json:"value"`
}

type PartialSimulation struct {
	SecondByDay            float64                        `json:"secondByDay"`
	ConsumptionsByHabitant []PartialConsumptionByHabitant `json:"consumptionsByHabitant"`
	Resources              []PartialResource              `json:"resources"`
	Countries              []PartialCountry               `json:"countries"`
	Territories            []PartialTerritory             `json:"territories"`
}

func (s *PartialSimulation) ToSimulation() Simulation {
	prices := make(map[ResourceType]float64, len(s.Resources))
	consumptionsByHabitant := make(map[ResourceType]float64, len(s.ConsumptionsByHabitant))
	for _, consumption := range s.ConsumptionsByHabitant {
		consumptionsByHabitant[consumption.Resource] = consumption.Amount
	}
	wg := new(sync.WaitGroup)
	for _, resource := range s.Resources {
		prices[resource.Name] = resource.Price
	}

	countries := make(map[string]*Country, len(s.Countries))
	for _, country := range s.Countries {
		in := make(Channel)
		out := make(Channel)
		//create slice of territories
		territories := make([]*Territory, 0)
		countries[country.ID] = NewCountry(country.ID, country.Name, country.Color, territories, country.Money, wg, in, out, consumptionsByHabitant)
	}

	territories := make([]*Territory, len(s.Territories))
	for i, territory := range s.Territories {
		var variations []Variation
		for _, variation := range territory.Variations {
			variations = append(variations, Variation{variation.Name, variation.Value})
		}
		country := countries[territory.Country]
		stock := make(map[ResourceType]float64)
		for _, stockT := range territory.Stocks {
			stock[stockT.Resource] = stockT.Amount
		}
		territories[i] = NewTerritory(territory.X, territory.Y, territory.Name, variations, stock, country, territory.Habitants)
		if country != nil {
			country.Territories = append(country.Territories, territories[i])
		}
	}
	return NewSimulation(s.SecondByDay, prices, countries, territories, wg, consumptionsByHabitant)
}

func LoadSimulation(path string) (Simulation, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Simulation{}, err
	}
	var simulation PartialSimulation
	err = json.Unmarshal(content, &simulation)
	if err != nil {
		return Simulation{}, err
	}
	return simulation.ToSimulation(), nil
}
