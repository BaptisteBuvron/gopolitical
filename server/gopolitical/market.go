package gopolitical

type Prices map[ResourceType]float64

type Market struct {
	sells  []MarketSellRequest
	buys   []MarketBuyRequest
	prices Prices
}

func NewMarket(prices Prices) Market {
	return Market{nil, nil, prices}
}

func (m *Market) handleRequest(req MarketBuyRequest) {
	// Add your implementation here
}
