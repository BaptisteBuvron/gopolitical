package gopolitical

import (
	"log"
	"math"
	"sync"
	"time"
)

type Event interface {
}

type CountryEvent struct {
	Event `json:"event"`
	Day   int `json:"day"`
}

type TransferResourceEvent struct {
	CountryEvent
	From     *Territory   `json:"from"`
	To       *Territory   `json:"to"`
	Resource ResourceType `json:"resource"`
	Amount   float64      `json:"amount"`
}

type Country struct {
	Agent        `json:"agent"`
	Color        string          `json:"color"`
	Territories  []*Territory    `json:"territories"`
	Money        float64         `json:"money"`
	History      []Event         `json:"history"`
	MoneyHistory map[int]float64 `json:"moneyHistory"`
	wg           *sync.WaitGroup `json:"-"`
	In           Channel         `json:"-"`
	Out          Channel         `json:"-"`
	currentDay   int             `json:"-"`
}

func NewCountry(id string, name string, color string, territories []*Territory, money float64, wg *sync.WaitGroup, in Channel, out Channel) *Country {
	return &Country{Agent{id, name}, color, territories, money, make([]Event, 0), make(map[int]float64), wg, in, out, 0}
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
	log.Printf("Country %s started\n", c.Name)
	for {
		c.wg.Add(1)
		c.Percept()
		requests := c.Deliberate()
		c.Act(requests)
		c.wg.Done()
		//Wait for the end of the day
		<-c.In
		//TODO: get the day from the environment (percept)
		c.currentDay++
	}
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

	log.Printf("Country %s percepted\n", c.Name)
}

func (c *Country) Deliberate() []Request {
	log.Printf("Country %s deliberate\n", c.Name)
	log.Println("Stock total de ", c.Name, " : ", c.GetTotalStock())
	time.Sleep(1 * time.Second)
	requests := []Request{}

	//Le pays regarde s'il lui manque des ressources, si oui, il les achète
	for _, territory := range c.Territories {
		foodConsomption := float64(territory.Habitants) * FOOD_BY_HABITANT
		waterConsumption := float64(territory.Habitants) * WATER_BY_HABITANT

		//Calculer si les territoires ont assez de ressources pour nourrir leurs habitants
		foodNeeded := territory.Stock["food"] - foodConsomption
		waterNeeded := territory.Stock["water"] - waterConsumption

		if foodNeeded < 0 {
			foodNeeded = math.Abs(foodNeeded)
			foodConsomption = c.tryTransferResources(territory, "food", foodNeeded)
		}
		if waterNeeded < 0 {
			waterNeeded = math.Abs(waterNeeded)
			waterConsumption = c.tryTransferResources(territory, "water", waterNeeded)
		}

		//TODO Rendre générique pour toutes les ressources
		if territory.Stock["food"] < foodConsomption {
			buyRequest := MarketBuyRequest{from: c, territoire: territory, resources: "food", amount: foodConsomption}
			log.Println("Ordre d'achat de ", foodConsomption, " food de ", c.Name, " pour le territoire ", territory.X, " ", territory.Y)
			requests = append(requests, buyRequest)
		}

		if territory.Stock["water"] < waterConsumption {
			buyRequest := MarketBuyRequest{from: c, territoire: territory, resources: "water", amount: waterConsumption}
			log.Println("Ordre d'achat de ", waterConsumption, " water de ", c.Name, " pour le territoire ", territory.X, " ", territory.Y)
			requests = append(requests, buyRequest)
		}

	}

	//Le pays regarde si des territoires ont plus de ressources que ce qu'il leur faut, si oui, il les vend
	for _, territory := range c.Territories {
		surplus := territory.GetSurplus(3)
		//Faire un ordre de vente pour chaque ressource en surplus
		for resource, quantity := range surplus {
			sellRequest := MarketSellRequest{from: c, territoire: territory, resources: resource, amount: quantity}
			log.Println("Ordre de vente de", quantity, " ", resource, " de ", c.Name, " pour le territoire ", territory.X, " ", territory.Y)
			requests = append(requests, sellRequest)
		}
	}

	return requests
}

func (c *Country) Act(requests []Request) {
	log.Printf("Country %s act\n", c.Name)
	for _, request := range requests {
		c.Out <- request
	}
}

func (c *Country) GetConsumption() map[ResourceType]float64 {
	consumption := make(map[ResourceType]float64)
	totalHabitants := c.GetTotalHabitants()

	consumption["foot"] = float64(totalHabitants) * FOOD_BY_HABITANT
	consumption["water"] = float64(totalHabitants) * WATER_BY_HABITANT

	return consumption
}

func (c *Country) tryTransferResources(to *Territory, resource ResourceType, need float64) float64 {
	for _, territory := range c.Territories {
		if territory != to {
			surplus := territory.GetSurplus(3)
			if surplus[resource] > 0 {
				if surplus[resource] > need {
					c.transferResources(territory, to, resource, need)
					return 0
				} else {
					c.transferResources(territory, to, resource, surplus[resource])
					return need - surplus[resource]
				}
			}
		}

	}
	return need
}

func (c *Country) transferResources(from *Territory, to *Territory, resource ResourceType, quantity float64) {
	log.Println("Transfert de ", quantity, " ", resource, " de ", from.Country.Name, " vers ", to.Country.Name)
	event := TransferResourceEvent{CountryEvent{"transferResource", c.currentDay}, from, to, resource, quantity}
	c.History = append(c.History, event)
	from.Stock[resource] -= quantity
	to.Stock[resource] += quantity
}
