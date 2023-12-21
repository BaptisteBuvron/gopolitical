package gopolitical

import (
	"encoding/json"
	"fmt"
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
		territories := make([]*Territory, 0)
		countries[country.ID] = NewCountry(country.ID, country.Name, country.Color, territories, country.Money, wg, in, out)
	}

	territories := make([]Territory, len(s.Territories))
	for i, territory := range s.Territories {
		var variations []Variation
		for _, variation := range territory.Variations {
			variations = append(variations, Variation{variation.Name, variation.Value})
		}
		country := countries[territory.Country]
		stock := make(map[ResourceType]int)
		for _, resource := range s.Resources {
			stock[resource.Name] = 0
		}
		territories[i] = NewTerritory(territory.X, territory.Y, variations, stock, country)
		country.Territories = append(country.Territories, &territories[i])
	}
	for _, country := range countries {
		fmt.Println("Nombre de territoires dans  : ", country.Name, " : ", len(country.Territories))
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
