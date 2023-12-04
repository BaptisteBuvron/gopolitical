package gopolitical

type Prices map[ResourceType]float64

type Market struct {
	sells  []MarketRessourceRequest
	buys   []MarketRessourceRequest
	prices Prices
}

func NewMarket(prices Prices) Market {
	return Market{nil, nil, prices}
}

type TypeRequest string

const (
	buyRequest TypeRequest = "buyRequest"
)

type Request struct {
	name TypeRequest
}

type MarketRessourceRequest struct {
	Request
	from      Country
	resources ResourceType
	amount    int
}
