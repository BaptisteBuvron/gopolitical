package gopolitical

import (
	"log"
)

type Prices map[ResourceType]float64

type Market struct {
	Env        *Environment `json:"-"`
	sells      map[int]*MarketSellRequest
	buys       map[int]*MarketBuyRequest
	Prices     Prices `json:"prices"`
	IdSell     int
	IdBuy      int
	Percept    map[string][]Request `json:"-"`
	History    []*MarketInteraction `json:"history"`
	currentDay int
}

func NewMarket(prices Prices, percept map[string][]Request) *Market {
	return &Market{nil, make(map[int]*MarketSellRequest), make(map[int]*MarketBuyRequest), prices, 0, 0, percept, []*MarketInteraction{}, 0}
}

func (m *Market) handleRequest(req MarketRequest) {
	switch req := req.(type) {
	case MarketBuyRequest:
		m.handleBuyRequest(&req)
	case MarketSellRequest:
		m.handleSellRequest(&req)
	default:
	}
}

func (m *Market) handleBuyRequest(req *MarketBuyRequest) {
	m.buys[m.IdBuy] = req
	req.BuyID = m.IdBuy
	m.IdBuy++
}

func (m *Market) handleSellRequest(req *MarketSellRequest) {
	m.sells[m.IdSell] = req
	req.SellID = m.IdSell
	m.IdSell++
}

func UpdateRelation(cost float64) {

}

func (m *Market) HandleRequests() {
	for _, buy := range m.buys {
		for _, sell := range m.sells {
			if buy.resources == sell.resources && buy.from != sell.from {
				cost := m.handleTransaction(buy, sell)
				if cost != 0.0 { // Optimization: Don't update if zero
					relation := m.Env.RelationManager.GetRelation(buy.from.ID, sell.from.ID)
					m.Env.RelationManager.UpdateRelation(buy.from.ID, sell.from.ID, relation+cost)
				}
			}
		}
	}
	log.Println(m.Env.RelationManager.ToString())
	//on vide les listes d'achats et de ventes pour le prochain tour => le pays recalculera au prochain tour ses ordres d'achats et de ventes en fonction de ses besoins
	m.buys = make(map[int]*MarketBuyRequest)
	m.sells = make(map[int]*MarketSellRequest)
}

func (m *Market) handleTransaction(buy *MarketBuyRequest, sell *MarketSellRequest) float64 {
	executed := 0.0
	if buy.amount >= sell.amount {
		executed = sell.amount
	} else {
		executed = buy.amount
	}

	//vérifier que le pays acheteur a assez d'argent
	if buy.from.Money < executed*m.Prices[buy.resources] {
		//on change la quantité executée
		executed = buy.from.Money / m.Prices[buy.resources]
		if executed <= 0 {
			log.Println("Transaction annulée : ", buy.from.Name, " n'a pas assez d'argent pour acheter ", executed, " ", buy.resources, " à ", sell.from.Name)
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

	//on envoie la reponse au pays acheteur
	m.Percept[buy.from.Name] = append(m.Percept[buy.from.Name], MarketBuyResponse{new(Request), "buyEvent", m.currentDay, buy.resources, sell.from.Name, executed, executed * m.Prices[buy.resources]})
	//on envoie la reponse au pays vendeur

	m.Percept[sell.from.Name] = append(m.Percept[sell.from.Name], MarketSellResponse{new(Request), "sellEvent", m.currentDay, sell.resources, buy.from.Name, executed, executed * m.Prices[sell.resources]})

	//on met à jour les stocks des pays et leur argent
	cost := executed * m.Prices[buy.resources]
	buy.from.Money -= cost
	sell.from.Money += cost

	buy.territoire.Stock[buy.resources] += executed
	sell.territoire.Stock[sell.resources] -= executed
	log.Println("Transaction effectuée : ", buy.from.Name, " achete ", executed, " ", buy.resources, " à ", sell.from.Name, " pour ", executed*m.Prices[buy.resources])
	m.History = append(m.History, &MarketInteraction{m.currentDay, buy.resources, executed, m.Prices[buy.resources], buy.from, sell.from})
	return cost
}
