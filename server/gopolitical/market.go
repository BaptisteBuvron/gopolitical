package gopolitical

import (
	"log"
	"math"
)

type Prices map[ResourceType]float64

type Market struct {
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
	return &Market{make(map[int]*MarketSellRequest), make(map[int]*MarketBuyRequest), prices, 0, 0, percept, []*MarketInteraction{}, 0}
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

func (m *Market) HandleRequests() {
	fusionHistory := make(map[string]map[string]map[ResourceType]*MarketInteraction)
	for _, buy := range m.buys {
		for _, sell := range m.sells {
			if buy.resources == sell.resources && buy.from != sell.from {
				m.handleTransaction(buy, sell, fusionHistory)
			}
		}
	}

	//add fusion history to market history
	for _, country := range fusionHistory {
		for _, country2 := range country {
			for _, marketInteraction := range country2 {
				m.History = append(m.History, marketInteraction)
			}
		}
	}

	//on vide les listes d'achats et de ventes pour le prochain tour => le pays recalculera au prochain tour ses ordres d'achats et de ventes en fonction de ses besoins
	m.buys = make(map[int]*MarketBuyRequest)
	m.sells = make(map[int]*MarketSellRequest)
}

func (m *Market) handleTransaction(buy *MarketBuyRequest, sell *MarketSellRequest, fusionHistory map[string]map[string]map[ResourceType]*MarketInteraction) {
	executed := 0.0
	if buy.amount >= sell.amount {
		executed = sell.amount
	} else {
		executed = buy.amount
	}

	//vérifier que le pays acheteur a assez d'argent
	if buy.from.Money < executed*m.Prices[buy.resources] {
		//on change la quantité executée
		executed = math.Floor(buy.from.Money / m.Prices[buy.resources])
		//get the integer part down
		if executed <= 0 {
			log.Println("Transaction annulée : ", buy.from.Name, " n'a pas assez d'argent pour acheter ", executed, " ", buy.resources, " à ", sell.from.Name)
			return
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

	//Add to history
	if fusionHistory[buy.from.Name] == nil {
		fusionHistory[buy.from.Name] = make(map[string]map[ResourceType]*MarketInteraction)
	}
	if fusionHistory[buy.from.Name][sell.from.Name] == nil {
		fusionHistory[buy.from.Name][sell.from.Name] = make(map[ResourceType]*MarketInteraction)
	}
	if fusionHistory[buy.from.Name][sell.from.Name][buy.resources] == nil {
		fusionHistory[buy.from.Name][sell.from.Name][buy.resources] = &MarketInteraction{m.currentDay, buy.resources, 0, m.Prices[buy.resources], buy.from, sell.from}
	}
	fusionHistory[buy.from.Name][sell.from.Name][buy.resources].Amount += executed

	//m.History = append(m.History, &MarketInteraction{m.currentDay, buy.resources, executed, m.Prices[buy.resources], buy.from, sell.from})
}
