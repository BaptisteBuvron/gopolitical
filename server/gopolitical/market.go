package gopolitical

type Prices map[ResourceType]float64

type Market struct {
	sells   map[int]MarketSellRequest
	buys    map[int]MarketBuyRequest
	prices  Prices
	idSell  int
	idBuy   int
	percept map[string][]Request
}

func NewMarket(prices Prices, percept map[string][]Request) Market {
	return Market{make(map[int]MarketSellRequest), make(map[int]MarketBuyRequest), prices, 0, 0, percept}
}

func (m *Market) handleRequest(req MarketRequest) {
	switch req := req.(type) {
	case MarketBuyRequest:
		m.handleBuyRequest(req)
	case MarketSellRequest:
		m.handleSellRequest(req)
	default:
	}
}

func (m *Market) handleBuyRequest(req MarketBuyRequest) {
	m.buys[m.idBuy] = req
	req.buyID = m.idBuy
	m.idBuy++
}

func (m *Market) handleSellRequest(req MarketSellRequest) {
	m.sells[m.idSell] = req
	req.sellID = m.idSell
	m.idSell++
}

func (m *Market) handleRequests() {
	for _, buy := range m.buys {
		for _, sell := range m.sells {
			if buy.resources == sell.resources {
				m.handleTransaction(buy, sell)
			}
		}
	}
}

func (m *Market) handleTransaction(buy MarketBuyRequest, sell MarketSellRequest) {
	executed := 0
	if buy.amount >= sell.amount {
		executed = sell.amount
	} else {
		executed = buy.amount
	}

	if executed == buy.amount {
		//Achat complet

		//on retire la vente de la liste des ventes
		sell.amount -= executed

		//on retire la demande d'achat de la liste des achats si elle est vide
		if sell.amount == 0 {
			delete(m.sells, sell.sellID)
		}
		//on retire la demande d'achat de la liste des achats
		delete(m.buys, buy.buyID)

		//on envoie la reponse au pays acheteur
		m.percept[buy.from.Name] = append(m.percept[buy.from.Name], MarketBuyResponse{buy, buy.from, executed, float64(executed) * m.prices[buy.resources]})

		//on envoie la reponse au pays vendeur
		m.percept[sell.from.Name] = append(m.percept[sell.from.Name], MarketSellResponse{sell, sell.from, executed, float64(executed) * m.prices[sell.resources]})

	} else {
		//Achat partiel

		//on modifie la vente
		sell.amount -= executed

		delete(m.sells, sell.sellID)

		//on modifie la demande d'achat
		buy.amount -= executed

		//on envoie la reponse au pays acheteur

		m.percept[buy.from.Name] = append(m.percept[buy.from.Name], MarketBuyResponse{buy, buy.from, executed, float64(executed) * m.prices[buy.resources]})

		//on envoie la reponse au pays vendeur
		m.percept[sell.from.Name] = append(m.percept[sell.from.Name], MarketSellResponse{sell, sell.from, executed, float64(executed) * m.prices[sell.resources]})
	}
}
