package gopolitical

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Country struct {
	Agent       `json:"agent"`
	Color       string          `json:"color"`
	Territories []*Territory    `json:"-"`
	Money       float64         `json:"money"`
	wg          *sync.WaitGroup `json:"-"`
	In          Channel         `json:"-"`
	Out         Channel         `json:"-"`
}

func NewCountry(id string, name string, color string, territories []*Territory, money float64, wg *sync.WaitGroup, in Channel, out Channel) *Country {
	return &Country{Agent{id, name}, color, territories, money, wg, in, out}
}

func (c *Country) GetTotalStock() map[ResourceType]float64 {
	stockCountry := make(map[ResourceType]float64)
	for _, territory := range c.Territories {
		for resource, quantity := range territory.Stock {
			stockCountry[resource] += quantity
		}
	}
	return stockCountry
}

func (c *Country) GetTotalHabitants() int {
	totalHabitants := 0
	for _, territory := range c.Territories {
		totalHabitants += territory.Habitants
	}
	return totalHabitants
}

func (c *Country) Start() {
	fmt.Printf("Country %s started\n", c.Name)
	c.Percept()
	requests := c.Deliberate()
	c.Act(requests)
	c.wg.Done()
}

func (c *Country) Percept() {
	perceptRequest := PerceptRequest{from: c}
	c.Out <- perceptRequest
	perceptResponse := <-c.In

	//Downcast to a PerceptResponse
	if perceptResponse, ok := perceptResponse.(PerceptResponse); ok {
		_ = perceptResponse
		//TODO : Faire un traitement des events
	}

	fmt.Printf("Country %s percepted\n", c.Name)
}

func (c *Country) Deliberate() []Request {
	fmt.Printf("Country %s deliberate\n", c.Name)
	fmt.Println("Stock total de ", c.Name, " : ", c.GetTotalStock())
	time.Sleep(1 * time.Second)

	//TODO : Faire des virements de ressources entre territoires du pays si besoin

	//Le pays regarde si des territoires ont plus de ressources que ce qu'il leur faut, si oui, il les vend
	for _, territory := range c.Territories {
		surplus := territory.GetSurplus()
		//Faire un ordre de vente pour chaque ressource en surplus
		for resource, quantity := range surplus {
			sellRequest := MarketSellRequest{from: c, territoire: territory, resources: resource, amount: quantity}
			log.Println("Ordre de vente de", quantity, " ", resource, " de ", c.Name, " pour le territoire ", territory.X, " ", territory.Y)
			//retirer les ressources du stock du territoire
			c.Out <- sellRequest
		}
	}

	//Le pays regarde s'il lui manque des ressources, si oui, il les achète
	for _, territory := range c.Territories {
		needFood := (float64(territory.Habitants) * FOOD_BY_HABITANT) - territory.Stock["food"]
		needWater := (float64(territory.Habitants) * WATER_BY_HABITANT) - territory.Stock["water"]

		if territory.Stock["food"] < float64(territory.Habitants)*FOOD_BY_HABITANT {
			buyRequest := MarketBuyRequest{from: c, territoire: territory, resources: "food", amount: needFood}
			log.Println("Ordre d'achat de ", needFood, " food de ", c.Name, " pour le territoire ", territory.X, " ", territory.Y)
			c.Out <- buyRequest
		}

		if territory.Stock["water"] < float64(territory.Habitants)*WATER_BY_HABITANT {
			buyRequest := MarketBuyRequest{from: c, territoire: territory, resources: "water", amount: needWater}
			log.Println("Ordre d'achat de ", needWater, " water de ", c.Name, " pour le territoire ", territory.X, " ", territory.Y)
			c.Out <- buyRequest
		}

	}

	return nil
}

func (c *Country) Act(requests []Request) {
	fmt.Printf("Country %s act\n", c.Name)
}

func (c *Country) GetConsumption() map[ResourceType]float64 {
	consumption := make(map[ResourceType]float64)
	totalHabitants := c.GetTotalHabitants()

	consumption["foot"] = float64(totalHabitants) * FOOD_BY_HABITANT
	consumption["water"] = float64(totalHabitants) * WATER_BY_HABITANT

	return consumption
}
