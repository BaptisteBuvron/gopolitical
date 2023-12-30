package gopolitical

import (
	"fmt"
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
	From     string       `json:"from"`
	To       string       `json:"to"`
	Resource ResourceType `json:"resource"`
	Amount   float64      `json:"amount"`
}

type BuyResourceEvent struct {
	CountryEvent
	From     string       `json:"from"`
	To       string       `json:"to"`
	Resource ResourceType `json:"resource"`
	Amount   float64      `json:"amount"`
	Price    float64      `json:"price"`
}

type SellResourceEvent struct {
	CountryEvent
	From     string       `json:"from"`
	To       string       `json:"to"`
	Resource ResourceType `json:"resource"`
	Amount   float64      `json:"amount"`
	Price    float64      `json:"price"`
}

type Country struct {
	Agent                 `json:"agent"`
	Color                 string                   `json:"color"`
	Territories           []*Territory             `json:"-"`
	Money                 float64                  `json:"money"`
	History               []Event                  `json:"history"`
	MoneyHistory          map[int]float64          `json:"moneyHistory"`
	wg                    *sync.WaitGroup          `json:"-"`
	In                    Channel                  `json:"-"`
	Out                   Channel                  `json:"-"`
	currentDay            int                      `json:"-"`
	consumptionByHabitant map[ResourceType]float64 `json:"-"`
}

func NewCountry(id string, name string, color string, territories []*Territory, money float64, wg *sync.WaitGroup, in Channel, out Channel, consumptionByHabitant map[ResourceType]float64) *Country {
	return &Country{Agent{id, name}, color, territories, money, make([]Event, 0), make(map[int]float64), wg, in, out, 0, consumptionByHabitant}
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
		//TODO: get the day from the environment (Percept)
		c.currentDay++
	}
}

func (c *Country) Percept() {
	perceptRequest := PerceptRequest{from: c}
	c.Out <- perceptRequest
	fmt.Println("Country ", c.Name, " waiting for percept")
	perceptResponse := <-c.In
	fmt.Println("Country ", c.Name, " percepted")

	//Downcast to a PerceptResponse
	if perceptResponse, ok := perceptResponse.(PerceptResponse); ok {
		//Process the events
		history := processMarketEvents(perceptResponse.events)
		c.History = append(c.History, history...)
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

		for resource, consumption := range c.consumptionByHabitant {
			totalConsumption := float64(territory.Habitants) * consumption
			//Calculer si les territoires ont assez de ressources pour nourrir leurs habitants
			needed := territory.Stock[resource] - totalConsumption
			if needed < 0 {
				needed = math.Abs(needed)
				consomption := c.tryTransferResources(territory, resource, needed)
				if consomption > 0 {
					buyRequest := MarketBuyRequest{from: c, territoire: territory, resources: resource, amount: consomption}
					log.Println("Ordre d'achat de ", consomption, " ", resource, " de ", c.Name, " pour le territoire ", territory.X, " ", territory.Y)
					requests = append(requests, buyRequest)
				}
			}
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
		//wait for the response to avoid concurrency problems
		<-c.In
	}
}

func (c *Country) tryTransferResources(to *Territory, resource ResourceType, need float64) float64 {
	for _, territory := range c.Territories {
		if territory != to {
			//Pour les échanges entre territoires, on ne prend que les surplus de 1 jour
			surplus := territory.GetSurplus(1)
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
	log.Println("Transfert de ", quantity, " ", resource, " de ", from.Name, " vers ", to.Name, "(", to.Country.Name, ")")
	event := TransferResourceEvent{CountryEvent{"transferResource", c.currentDay}, from.Name, to.Name, resource, quantity}
	c.History = append(c.History, event)
	from.Stock[resource] -= quantity
	to.Stock[resource] += quantity
}

func processMarketEvents(events []Request) []Event {
	buyEvents := make(map[string]map[ResourceType]MarketBuyResponse)
	sellEvents := make(map[string]map[ResourceType]MarketSellResponse)

	for _, event := range events {
		switch event := event.(type) {
		case MarketBuyResponse:
			processBuyEvent(event, buyEvents)
		case MarketSellResponse:
			processSellEvent(event, sellEvents)
		}
	}

	return combineEvents(buyEvents, sellEvents)
}

func processBuyEvent(event MarketBuyResponse, buyEvents map[string]map[ResourceType]MarketBuyResponse) {
	if _, ok := buyEvents[event.From]; !ok {
		buyEvents[event.From] = make(map[ResourceType]MarketBuyResponse)
	}

	if _, ok := buyEvents[event.From][event.ResourceType]; !ok {
		buyEvents[event.From][event.ResourceType] = event
	} else {
		ev := buyEvents[event.From][event.ResourceType]
		ev.AmountExecuted += event.AmountExecuted
		ev.Cost += event.Cost
		buyEvents[event.From][event.ResourceType] = ev
	}
}

func processSellEvent(event MarketSellResponse, sellEvents map[string]map[ResourceType]MarketSellResponse) {
	if _, ok := sellEvents[event.To]; !ok {
		sellEvents[event.To] = make(map[ResourceType]MarketSellResponse)
	}

	if _, ok := sellEvents[event.To][event.ResourceType]; !ok {
		sellEvents[event.To][event.ResourceType] = event
	} else {
		ev := sellEvents[event.To][event.ResourceType]
		ev.AmountExecuted += event.AmountExecuted
		ev.Gain += event.Gain
		sellEvents[event.To][event.ResourceType] = ev
	}
}

func combineEvents(buyEvents map[string]map[ResourceType]MarketBuyResponse, sellEvents map[string]map[ResourceType]MarketSellResponse) []Event {
	var combinedEvents []Event

	for _, country := range buyEvents {
		for _, event := range country {
			combinedEvents = append(combinedEvents, event)
		}
	}

	for _, country := range sellEvents {
		for _, event := range country {
			combinedEvents = append(combinedEvents, event)
		}
	}

	return combinedEvents
}
