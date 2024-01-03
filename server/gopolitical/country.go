package gopolitical

import (
	"fmt"
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
	Flag                  string                   `json:"flag"`
	Territories           []*Territory             `json:"-"`
	Money                 float64                  `json:"money"`
	History               []Event                  `json:"history"`
	MoneyHistory          map[int]float64          `json:"moneyHistory"`
	RelationManager       *RelationManager         `json:"-"`
	PerceivedWorld        *World                   `json:"-"`
	PerceivedPrices       Prices                   `json:"-"`
	In                    Channel                  `json:"-"`
	Out                   Channel                  `json:"-"`
	CurrentDay            int                      `json:"-"`
	WgStart               *sync.WaitGroup          `json:"-"`
	WgMiddle              *sync.WaitGroup          `json:"-"`
	WgEnd                 *sync.WaitGroup          `json:"-"`
	RandomGenerator       *rand.Rand               `json:"-"`
	consumptionByHabitant map[ResourceType]float64 `json:"-"`
}

func NewCountry(id string, name string, color string, flag string, territories []*Territory, money float64, in Channel, out Channel, consumptionByHabitant map[ResourceType]float64) *Country {
	country := &Country{
		Agent:                 Agent{id, name},
		Color:                 color,
		Flag:                  flag,
		Territories:           territories,
		Money:                 money,
		History:               make([]Event, 0),
		RelationManager:       nil,
		MoneyHistory:          make(map[int]float64),
		PerceivedWorld:        nil,
		PerceivedPrices:       nil,
		In:                    in,
		Out:                   out,
		CurrentDay:            0,
		consumptionByHabitant: consumptionByHabitant,
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
	Debug(c.Name, "Lancée")
	for {
		c.WgStart.Wait()
		Debug(c.Name, "Commence ses actions")

		c.Percept()
		requests := c.Deliberate()
		c.Act(requests)
		//Wait for the end of the day
		//TODO: get the day from the environment (Percept)
		c.CurrentDay++

		Debug(c.Name, "Termine ses actions")
		c.WgMiddle.Done()
		c.WgEnd.Wait()
		c.WgMiddle.Done()
		Debug(c.Name, "Attend le prochain jour")
	}
}

func (c *Country) Percept() {
	perceptRequest := PerceptRequest{from: c}
	c.Out <- perceptRequest
	Debug(c.Name, "Percept")
	perceptResponse := <-c.In
	//Downcast to a PerceptResponse
	if perceptResponse, ok := perceptResponse.(PerceptResponse); ok {
		//Process the events
		history := processMarketEvents(perceptResponse.events)
		c.History = append(c.History, history...)
		c.RelationManager = perceptResponse.RelationManager
		c.PerceivedWorld = perceptResponse.World
		c.PerceivedPrices = perceptResponse.Prices
	} else {
		// TODO: handle error
	}
}

func (c *Country) Deliberate() []Request {
	Debug(c.Name, "Deliberate")
	Debug(c.Name, "Stock %v", c.GetTotalStock())
	time.Sleep(1 * time.Second)
	requests := []Request{}

	//Le pays regarde s'il lui manque des ressources, si oui, il les achète
	for _, territory := range c.Territories {
		for resource, consumption := range c.consumptionByHabitant {
			totalConsumption := (float64(territory.Habitants) * consumption) * 2
			//Calculer si les territoires ont assez de ressources pour nourrir leurs habitants
			needed := territory.Stock[resource] - totalConsumption
			if needed < 0 {
				needed = math.Abs(needed)
				consomption := c.tryTransferResources(territory, resource, needed)
				if consomption > 0 {
					buyRequest := MarketBuyRequest{from: c, territoire: territory, resources: resource, amount: consomption}
					Debug(c.Name, "Ordre d'achat de %f %s via %s", consomption, resource, territory.Name)
					requests = append(requests, buyRequest)
				}
			}
		}
		// On achète de l’armement pour 10 par pays
		armement := c.GetTotalStockOf(ARMAMENT)
		armementRequired := float64(len(c.Territories) * ARMAMENT_NEEDED_BY_TERRITORY)
		if armement < armementRequired {
			buyRequest := MarketBuyRequest{from: c, territoire: territory, resources: "armement", amount: armementRequired - armement}
			Debug(c.Name, "Ordre d'achat de %f armement via %s", armementRequired-armement, territory.Name)
			requests = append(requests, buyRequest)
		} else if armement > armementRequired {
			sellRequest := MarketSellRequest{from: c, territoire: territory, resources: "armement", amount: armement - armementRequired}
			Debug(c.Name, "Ordre de vente de %f armement via %s", armement-armementRequired, territory.Name)
			requests = append(requests, sellRequest)
		}
	}

	//Le pays regarde si des territoires ont plus de ressources que ce qu'il leur faut, si oui, il les vend
	for _, territory := range c.Territories {
		surplus := territory.GetSurplus(DAYS_TO_SECURE)
		//Faire un ordre de vente pour chaque ressource en surplus
		for resource, quantity := range surplus {
			sellRequest := MarketSellRequest{from: c, territoire: territory, resources: resource, amount: quantity}
			Debug(c.Name, "Ordre de vente de %f %s via %s", quantity, resource, territory.Name)
			requests = append(requests, sellRequest)
		}
	}

	// Check for war >:[
	stock := c.GetTotalStock()
	for resource, quantity := range stock {
		resourceConsumption := c.consumptionByHabitant[resource]
		missing := quantity - resourceConsumption*DAYS_TO_WARS
		if missing < 0 && stock[ARMAMENT] > 0 {
			territory := c.MostInterestingTerritoryToAttack()
			if territory != nil {
				attackReq := AttackRequest{from: c, to: territory, armement: c.RandomGenerator.Float64() * stock[ARMAMENT]}
				requests = append(requests, attackReq)
				Debug(c.Name, "pense à attaquer %v {%d, %d} avec %.2f armement", territory.Country.ID, territory.X, territory.Y, attackReq.armement)
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
	Debug(c.Name, "Act")
	for _, request := range requests {
		c.Out <- request
		//wait for the response to avoid concurrency problems
		<-c.In
	}
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
			//Pour les échanges entre territoires, on ne prend que les surplus de 1 jour
			surplus := territory.GetSurplus(2)
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
	Debug(from.Name, "Transfert de %f %s vers %s (%s) ", quantity, resource, to.Name, to.Country.Name)
	event := TransferResourceEvent{CountryEvent{"transferResource", c.CurrentDay}, from.Name, to.Name, resource, quantity}
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
