package gopolitical

import (
	"math"
)

type Prices map[ResourceType]float64

type MarketBuyAction struct {
	Action
	BuyID      int
	from       *Country
	territoire *Territory
	resources  ResourceType
	amount     float64
}

func (a *MarketBuyAction) Execute(env *Environment) {
	env.Market.RegisterBuy(a)
}

type MarketSellAction struct {
	Action
	SellID     int
	from       *Country
	territoire *Territory
	resources  ResourceType
	amount     float64
}

func (a *MarketSellAction) Execute(env *Environment) {
	env.Market.RegisterSell(a)
}

type Market struct {
	Env        *Environment `json:"-"`
	sells      map[int]*MarketSellAction
	buys       map[int]*MarketBuyAction
	Prices     Prices `json:"prices"`
	IdSell     int
	IdBuy      int
	History    []*MarketInteractionEvent `json:"history"`
	currentDay int
}

func (m *Market) RegisterBuy(buy *MarketBuyAction) {
	m.buys[m.IdBuy] = buy
	buy.BuyID = m.IdSell
	m.IdSell++
}

func (m *Market) RegisterSell(sell *MarketSellAction) {
	m.sells[m.IdSell] = sell
	sell.SellID = m.IdSell
	m.IdSell++
}
func NewMarket(prices Prices) *Market {
	return &Market{
		nil,
		make(map[int]*MarketSellAction),
		make(map[int]*MarketBuyAction),
		prices,
		0,
		0,
		[]*MarketInteractionEvent{},
		0,
	}
}

func (m *Market) HandleRequests() {
	fusionHistory := make(map[string]map[string]map[ResourceType]*MarketInteractionEvent)
	for _, buy := range m.buys {
		for _, sell := range m.sells {
			if buy.resources == sell.resources && buy.from != sell.from {
				cost := m.HandleTransaction(buy, sell, fusionHistory)
				if cost != 0.0 { // Optimization: Don't update if zero
					relation := m.Env.RelationManager.GetRelation(buy.from.ID, sell.from.ID)
					m.Env.RelationManager.UpdateRelation(buy.from.ID, sell.from.ID, relation+cost)
				}
			}
		}
	}

	//remove history higher than 4 days
	newHistory := make([]*MarketInteractionEvent, 0)
	for _, country := range m.History {
		if country.DateTransaction >= m.currentDay-4 {
			newHistory = append(newHistory, country)
		}
	}
	m.History = newHistory

	//add fusion history to market history
	for _, country := range fusionHistory {
		for _, country2 := range country {
			for _, marketInteraction := range country2 {
				m.History = append(m.History, marketInteraction)
			}
		}
	}
	//on vide les listes d'achats et de ventes pour le prochain tour => le pays recalculera au prochain tour ses ordres d'achats et de ventes en fonction de ses besoins
	m.buys = make(map[int]*MarketBuyAction)
	m.sells = make(map[int]*MarketSellAction)
}

func (m *Market) HandleTransaction(
	buy *MarketBuyAction,
	sell *MarketSellAction,
	fusionHistory map[string]map[string]map[ResourceType]*MarketInteractionEvent,
) float64 {
	executed := 0.0
	if buy.amount >= sell.amount {
		executed = sell.amount
	} else {
		executed = buy.amount
	}
	// On vérifie que le territoire n'a pas changé
	if buy.from != buy.territoire.Country || sell.from != sell.territoire.Country {
		Debug("Market", "[%s->%s] Transaction invalide de %f %s", buy.from.Name, sell.from.Name, executed, buy.resources)
		return 0.0

	}

	//vérifier que le pays acheteur a assez d'argent
	if buy.from.Money < executed*m.Prices[buy.resources] {
		//on change la quantité executée
		executed = math.Floor(buy.from.Money / m.Prices[buy.resources])
		//get the integer part down
		if executed <= 0 {
			Debug("Market", "[%s->%s] Transaction annulée de %f %s", buy.from.Name, sell.from.Name, executed, buy.resources)
			return 0.0
		}
	}

	if executed == buy.amount {
		//Achat complet

		//on retire la vente de la liste des ventes
		sell.amount -= executed

		//on retire la demande d'achat de la liste des achats si elle est vide
		if sell.amount == 0 {
			delete(m.sells, sell.SellID)
		}
		//on retire la demande d'achat de la liste des achats
		delete(m.buys, buy.BuyID)

	} else {
		//Achat partiel

		//on modifie la vente
		sell.amount -= executed

		delete(m.sells, sell.SellID)

		//on modifie la demande d'achat
		buy.amount -= executed

	}

	//on met à jour les stocks des pays et leur argent
	cost := executed * m.Prices[buy.resources]
	buy.from.Money -= cost
	sell.from.Money += cost

	buy.territoire.Stock[buy.resources] += executed
	sell.territoire.Stock[sell.resources] -= executed
	Debug("Market", "[%s->%s] Transaction effectuée de %f %s pour %f", buy.from.Name, sell.from.Name, executed, buy.resources, cost)
	//D: m.History = append(m.History, &MarketInteraction{m.currentDay, buy.resources, executed, m.Prices[buy.resources], buy.from, sell.from})

	buyer := buy.from.Name
	seller := sell.from.Name
	//Add to history
	if fusionHistory[buyer] == nil {
		fusionHistory[buyer] = make(map[string]map[ResourceType]*MarketInteractionEvent)
	}
	if fusionHistory[buyer][seller] == nil {
		fusionHistory[buyer][seller] = make(map[ResourceType]*MarketInteractionEvent)
	}
	if fusionHistory[buyer][seller][buy.resources] == nil {
		record := &MarketInteractionEvent{m.currentDay, buy.resources, 0, m.Prices[buy.resources], buy.from, sell.from}
		fusionHistory[buyer][seller][buy.resources] = record
	}
	fusionHistory[buyer][seller][buy.resources].Amount += executed

	return cost
}
