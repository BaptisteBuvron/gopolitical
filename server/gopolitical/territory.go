package gopolitical

import "log"

type Territory struct {
	X                int                              `json:"x"`
	Y                int                              `json:"y"`
	Name             string                           `json:"name"`
	Variations       []Variation                      `json:"variations"`
	Stock            map[ResourceType]float64         `json:"stock"`
	StockHistory     map[int]map[ResourceType]float64 `json:"stockHistory"`
	Country          *Country                         `json:"country"`
	Habitants        int                              `json:"habitants"`
	HabitantsHistory map[int]int                      `json:"habitantsHistory"`
}

func NewTerritory(x int, y int, name string, variations []Variation, stock map[ResourceType]float64, country *Country, habitant int) *Territory {
	return &Territory{x, y, name, variations, stock, make(map[int]map[ResourceType]float64), country, habitant, make(map[int]int)}
}

func (t Territory) Start() {

}

func (t Territory) MarketValue(prices Prices) float64 {
	totalVariation := 0.0
	for _, variation := range t.Variations {
		totalVariation += variation.Amount * prices[variation.Ressource]
	}
	return totalVariation
}

func (t *Territory) GetSurplus(daysToSecure float64) map[ResourceType]float64 {
	surplus := make(map[ResourceType]float64)
	//On garde 3 jours de surplus

	for resource, consumption := range t.Country.consumptionByHabitant {
		surplusAmount := t.Stock[resource] - (float64(t.Habitants)*consumption)*daysToSecure
		if surplusAmount > 0 {
			surplus[resource] = surplusAmount
		}
	}
	return surplus
}

func (t *Territory) TransfertProperty(country *Country) {
	log.Printf("[%s] Transfert de propriete de %s vers %s", t.Name, t.Country.Name, country.Name)
	losingCountry := t.Country
	// Délier le territoire a son pays d'origine
	for i, territory := range losingCountry.Territories {
		if territory.Equal(t) {
			losingCountry.Territories = append(losingCountry.Territories[:i], losingCountry.Territories[i+1:]...)
		}
	}
	t.Country = country
	// Relier le territoire au nouveaux pays
	country.Territories = append(country.Territories, t)

	// On enlève les dettes du territoires et les faisons consommer au pays
	for ressource, quantity := range t.Stock {
		if quantity < 0 {
			losingCountry.Consume(ressource, -quantity)
		}
	}
	// TODO: migrate population
	t.Habitants = 0
}
func (t *Territory) Equal(territory *Territory) bool {
	return t.X == territory.X && t.Y == territory.Y
}
