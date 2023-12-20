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
	X          int                `json:"id"`
	Y          int                `json:"color"`
	Country    string             `json:"country"`
	Variations []PartialVariation `json:"variations"`
}

type PartialVariation struct {
	Name  ResourceType `json:"name"`
	Value float64      `json:"value"`
}

type PartialSimulation struct {
	SecondByDay float64            `json:"secondByDay"`
	Resources   []PartialResource  `json:"resources"`
	Countries   []PartialCountry   `json:"countries"`
	Territories []PartialTerritory `json:"territories"`
}

func (s *PartialSimulation) ToSimulation() Simulation {
	prices := make(map[ResourceType]float64, len(s.Resources))
	wg := new(sync.WaitGroup)
	for _, resource := range s.Resources {
		prices[resource.Name] = resource.Price
	}

	countries := make(map[string]Country, len(s.Countries))
	for _, country := range s.Countries {
		in := make(Channel)
		out := make(Channel)
		countries[country.ID] = NewCountry(country.ID, country.Name, country.Color, nil, country.Money, wg, in, out)
	}

	territories := make([]Territory, len(s.Territories))
	for i, territory := range s.Territories {
		var variations []Variation
		for _, variation := range territory.Variations {
			variations = append(variations, Variation{variation.Name, variation.Value})
		}
		country := countries[territory.Country]
		territories[i] = NewTerritory(territory.X, territory.Y, variations, country)
		country.Territories = append(country.Territories, territories[i])
	}

	return NewSimulation(s.SecondByDay, prices, countries, territories, wg)
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
