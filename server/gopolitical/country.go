package gopolitical

import (
	"fmt"
	"log"
	"math"
	"math/rand"
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

type BuyResourceEvent struct {
	CountryEvent
	From     *Territory   `json:"from"`
	To       *Territory   `json:"to"`
	Resource ResourceType `json:"resource"`
	Amount   float64      `json:"amount"`
	Price    float64      `json:"price"`
}

type SellResourceEvent struct {
	CountryEvent
	From     *Territory   `json:"from"`
	To       *Territory   `json:"to"`
	Resource ResourceType `json:"resource"`
	Amount   float64      `json:"amount"`
	Price    float64      `json:"price"`
}

type Country struct {
	Agent           `json:"agent"`
	Color           string           `json:"color"`
	Territories     []*Territory     `json:"territories"`
	Money           float64          `json:"money"`
	History         []Event          `json:"history"`
	MoneyHistory    map[int]float64  `json:"moneyHistory"`
	RelationManager *RelationManager `json:"-"`
	PerceivedWorld  *World           `json:"-"`
	PerceivedPrices Prices           `json:"-"`
	In              Channel          `json:"-"`
	Out             Channel          `json:"-"`
	CurrentDay      int              `json:"-"`
	WgStart         *sync.WaitGroup  `json:"-"`
	WgMiddle        *sync.WaitGroup  `json:"-"`
	WgEnd           *sync.WaitGroup  `json:"-"`
	RandomGenerator *rand.Rand       `json:"-"`
}

func NewCountry(id string, name string, color string, territories []*Territory, money float64, in Channel, out Channel) *Country {
	country := &Country{
		Agent:           Agent{id, name},
		Color:           color,
		Territories:     territories,
		Money:           money,
		History:         make([]Event, 0),
		RelationManager: nil,
		MoneyHistory:    make(map[int]float64),
		PerceivedWorld:  nil,
		PerceivedPrices: nil,
		In:              in,
		Out:             out,
		CurrentDay:      0,
	}
	randomSource := rand.NewSource(time.Now().UnixNano())
	country.RandomGenerator = rand.New(randomSource)
	return country
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

func (c *Country) GetTotalStockOf(ressource ResourceType) float64 {
	stock := 0.0
	for _, territory := range c.Territories {
		stock += territory.Stock[ressource]
	}
	return stock
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
		c.WgStart.Wait()
		log.Printf("Start %s\n", c.Name)

		c.Percept()
		requests := c.Deliberate()
		c.Act(requests)
		//Wait for the end of the day
		//TODO: get the day from the environment (Percept)
		c.CurrentDay++

		log.Printf("End %s\n", c.Name)
		c.WgMiddle.Done()
		c.WgEnd.Wait()
		c.WgMiddle.Done()
		log.Printf("Wait next day %s\n", c.Name)
	}
}

func (c *Country) Percept() {
	perceptRequest := PerceptRequest{from: c}
	c.Out <- perceptRequest
	perceptResponse := <-c.In

	//Downcast to a PerceptResponse
	if perceptResponse, ok := perceptResponse.(PerceptResponse); ok {
		for _, event := range perceptResponse.events {
			switch event.(type) {
			case MarketBuyResponse, MarketSellResponse:
				c.History = append(c.History, event)
				break
			}
		}
		c.RelationManager = perceptResponse.RelationManager
		c.PerceivedWorld = perceptResponse.World
		c.PerceivedPrices = perceptResponse.Prices
		//TODO : Faire un traitement des events
	} else {
		panic("Invalid state, got unexpected response for percept request")
	}

	log.Printf("Country %s perceived\n", c.Name)
}

func (c *Country) Deliberate() []Request {
	log.Printf("Country %s deliberate\n", c.Name)
	log.Println("Stock total de", c.Name, ":", c.GetTotalStock())
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
			log.Println("Ordre d'achat de", foodConsomption, "food de", c.Name, " our le territoire", territory.X, " ", territory.Y)
			requests = append(requests, buyRequest)
		}

		if territory.Stock["water"] < waterConsumption {
			buyRequest := MarketBuyRequest{from: c, territoire: territory, resources: "water", amount: waterConsumption}
			log.Println("Ordre d'achat de", waterConsumption, "water de", c.Name, "pour le territoire", territory.X, " ", territory.Y)
			requests = append(requests, buyRequest)
		}
	}

	//Le pays regarde si des territoires ont plus de ressources que ce qu'il leur faut, si oui, il les vend
	for _, territory := range c.Territories {
		surplus := territory.GetSurplus(3)
		//Faire un ordre de vente pour chaque ressource en surplus
		for resource, quantity := range surplus {
			sellRequest := MarketSellRequest{from: c, territoire: territory, resources: resource, amount: quantity}
			log.Println("Ordre de vente de", quantity, resource, "de", c.Name, "pour le territoire", territory.X, " ", territory.Y)
			requests = append(requests, sellRequest)
		}
	}

	// Check for war >:[
	stock := c.GetTotalStock()
	consumption := c.GetConsumption()
	facteur := 1.0
	for resource, quantity := range stock {
		resourceConsumption := consumption[resource]
		missing := quantity - resourceConsumption*facteur
		if missing < 0 {
			territory := c.MostInterestingTerritoryToAttack()
			if territory != nil {
				attackReq := AttackRequest{from: c, to: territory, armement: c.RandomGenerator.Float64() * stock[ARMAMENT]}
				requests = append(requests, attackReq)
				log.Printf("%v pense à attaquer %v %v %v\n", c.ID, territory.X, territory.Y, territory.Country.ID)
			}
			break
		}
	}

	return requests
}

func (c *Country) MostInterestingTerritoryToAttack() *Territory {
	// Trouve un territoire voisin avec le moins de défense et le plus de ressources
	var bestAttackTerritory *Territory
	bestAttackScore := 0.0
	for _, territory := range c.PerceivedWorld.FindNeighborTerritoriesOfCountry(c) {
		relation := c.RelationManager.GetRelation(c.ID, territory.Country.ID)
		value := territory.MarketValue(c.PerceivedPrices)
		stock := territory.Country.GetTotalStock()
		attackScore := (1 / relation) * value * (1 / stock[ARMAMENT])
		if attackScore > bestAttackScore {
			bestAttackTerritory = territory
			bestAttackScore = attackScore
		}
	}
	return bestAttackTerritory
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

// O(3 * Territories)
func (c *Country) Consume(resource ResourceType, quantity float64) error {
	if len(c.Territories) == 0 {
		// On ignore la demande
		return fmt.Errorf("Aucuns territoire valide")
	}
	// On cherche le stock minimum (qui peut être négatif)
	minStock := c.Territories[0].Stock[resource]
	for _, territory := range c.Territories[1:] {
		stock := territory.Stock[resource]
		if minStock > stock {
			minStock = stock
		}
	}
	// On calcule la somme des différences avec le plus petit
	totalStockDifference := 0.0
	for _, territory := range c.Territories {
		stock := territory.Stock[resource]
		totalStockDifference += stock - minStock
	}

	// On consomme les ressources en fonction du ratio stock / totalStockDifference
	for _, territory := range c.Territories {
		stock := territory.Stock[resource]
		var percent float64
		if totalStockDifference == 0.0 {
			percent = 1.0 / float64(len(c.Territories))
		} else {
			percent = (stock - minStock) / totalStockDifference
		}
		territory.Stock[resource] -= quantity * percent
	}
	return nil
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
	log.Println("Transfert de", quantity, resource, "de", from.Country.Name, "vers", to.Country.Name)
	event := TransferResourceEvent{CountryEvent{"transferResource", c.CurrentDay}, from, to, resource, quantity}
	c.History = append(c.History, event)
	from.Stock[resource] -= quantity
	to.Stock[resource] += quantity
}
