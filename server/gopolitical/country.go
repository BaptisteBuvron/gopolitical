package gopolitical

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Country struct {
	Agent
	ID              string                   `json:"id"`
	Name            string                   `json:"name"`
	Color           string                   `json:"color"`
	Flag            string                   `json:"flag"`
	Territories     []*Territory             `json:"-"`
	Money           float64                  `json:"money"`
	History         []Event                  `json:"history"`
	MoneyHistory    map[int]float64          `json:"moneyHistory"`
	Perception      *CompletePerception      `json:"-"`
	In              chan *CompletePerception `json:"-"`
	Out             chan []Action            `json:"-"`
	RandomGenerator *rand.Rand               `json:"-"`
}

func NewCountry(
	id string,
	name string,
	color string,
	flag string,
	territories []*Territory,
	money float64,
	in chan *CompletePerception,
	out chan []Action,
	consumptionByHabitant map[ResourceType]float64,
) *Country {
	country := &Country{
		ID:           id,
		Name:         name,
		Color:        color,
		Flag:         flag,
		Territories:  territories,
		Money:        money,
		History:      make([]Event, 0),
		MoneyHistory: make(map[int]float64),
		In:           in,
		Out:          out,
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
		// Get from Percept()
		c.Perception = <-c.In
		Debug(c.Name, "Commence ses actions")

		//remove history higher than 4 days
		newHistory := make([]Event, 0)
		for _, country := range c.History {
			if country.(TransferResourceEvent).Day >= c.Perception.env.CurrentDay-4 {
				newHistory = append(newHistory, country)
			}
		}
		c.History = newHistory

		// Deliberate will call <- Out
		c.Deliberate()
	}
}

func (c *Country) Percept(env *Environment) {
	Debug(c.Name, "Percept")
	c.In <- &CompletePerception{env}
}

func (c *Country) Deliberate() {
	Debug(c.Name, "Deliberate")
	Debug(c.Name, "Stock %v", c.GetTotalStock())
	time.Sleep(1 * time.Second)
	actions := []Action{}

	//Le pays regarde s'il lui manque des ressources, si oui, il les achète
	for _, territory := range c.Territories {
		for resource, consumption := range c.Perception.env.ConsumptionByHabitant {
			totalConsumption := (float64(territory.Habitants) * consumption) * 2
			//Calculer si les territoires ont assez de ressources pour nourrir leurs habitants
			needed := territory.Stock[resource] - totalConsumption
			if needed < 0 {
				needed = math.Abs(needed)
				consomption := c.tryTransferResources(territory, resource, needed)
				if consomption > 0 {
					buy := &MarketBuyAction{from: c, territoire: territory, resources: resource, amount: consomption}
					Debug(c.Name, "Ordre d'achat de %f %s via %s", consomption, resource, territory.Name)
					actions = append(actions, buy)
				}
			}
		}
		// On achète de l’armement pour 10 par pays
		armement := c.GetTotalStockOf(ARMAMENT)
		armementRequired := float64(len(c.Territories) * ARMAMENT_NEEDED_BY_TERRITORY)
		if armement < armementRequired {
			buy := &MarketBuyAction{from: c, territoire: territory, resources: "armement", amount: armementRequired - armement}
			Debug(c.Name, "Ordre d'achat de %f armement via %s", armementRequired-armement, territory.Name)
			actions = append(actions, buy)
		} else if armement > armementRequired {
			sell := &MarketSellAction{from: c, territoire: territory, resources: "armement", amount: armement - armementRequired}
			Debug(c.Name, "Ordre de vente de %f armement via %s", armement-armementRequired, territory.Name)
			actions = append(actions, sell)
		}
	}

	//Le pays regarde si des territoires ont plus de ressources que ce qu'il leur faut, si oui, il les vend
	for _, territory := range c.Territories {
		surplus := territory.GetSurplus(DAYS_TO_SECURE, c.Perception.env.ConsumptionByHabitant)
		//Faire un ordre de vente pour chaque ressource en surplus
		for resource, quantity := range surplus {
			sell := &MarketSellAction{from: c, territoire: territory, resources: resource, amount: quantity}
			Debug(c.Name, "Ordre de vente de %f %s via %s", quantity, resource, territory.Name)
			actions = append(actions, sell)
		}
	}

	// Check for war >:[
	stock := c.GetTotalStock()
	for resource, quantity := range stock {
		resourceConsumption := c.Perception.env.ConsumptionByHabitant[resource]
		missing := quantity - resourceConsumption*DAYS_TO_WARS
		if missing < 0 && stock[ARMAMENT] > 0 {
			territory := c.MostInterestingTerritoryToAttack()
			if territory != nil {
				attack := &AttackAction{
					AttackerCountry:   c,
					DefenderTerritory: territory,
					ArmamentUsed:      c.RandomGenerator.Float64() * stock[ARMAMENT],
				}
				actions = append(actions, attack)
				Debug(c.Name, "pense à attaquer %v {%d, %d} avec %.2f armement", territory.Country.ID, territory.X, territory.Y, attack.ArmamentUsed)
			}
			break
		}
	}

	c.Out <- actions
}

func (c *Country) MostInterestingTerritoryToAttack() *Territory {
	// Trouve un territoire voisin avec le moins de défense et le plus de ressources
	var bestAttackTerritory *Territory
	bestAttackScore := math.Inf(-1)
	for _, territory := range c.Perception.env.World.FindNeighborTerritoriesOfCountry(c) {
		relation := c.Perception.env.RelationManager.GetRelation(c.ID, territory.Country.ID)
		value := territory.MarketValue(c.Perception.env.Market.Prices)
		armament := 1.0 // territory.Country.GetTotalStockOf(ARMAMENT)
		bonusContact := (1.0 + float64(c.Perception.env.World.CountDirectContact(territory, c)))
		attackScore := (1.0 / (100.0 + relation)) * value * (1.0 / (100.0 + armament)) * bonusContact
		if attackScore > bestAttackScore {
			bestAttackTerritory = territory
			bestAttackScore = attackScore
		}
	}
	return bestAttackTerritory
}

func (c *Country) Act() []Action {
	Debug(c.Name, "Act")
	return <-c.Out
}

func (c *Country) GetID() string {
	return c.ID
}

func (c *Country) CleanUp(env *Environment) {
	// TODO
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
			surplus := territory.GetSurplus(2, c.Perception.env.ConsumptionByHabitant)
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
	event := TransferResourceEvent{CountryEvent{"transferResource", c.Perception.env.CurrentDay}, from.Name, to.Name, resource, quantity}
	c.History = append(c.History, event)
	from.Stock[resource] -= quantity
	to.Stock[resource] += quantity
}
