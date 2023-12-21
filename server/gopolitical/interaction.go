package gopolitical

type Channel = chan Request

type Request interface {
}

//Market

type MarketRequest interface {
	Request
}

type MarketBuyRequest struct {
	MarketRequest
	buyID     int
	from      Country
	resources ResourceType
	amount    int
}

type MarketSellRequest struct {
	MarketRequest
	sellID    int
	from      Country
	resources ResourceType
	amount    int
}

type MarketBuyResponse struct {
	MarketRequest
	to             Country
	amountExecuted int
	cost           float64
}

type MarketSellResponse struct {
	MarketRequest
	to             Country
	amountExecuted int
	gain           float64
}

//Percept

type PerceptRequest struct {
	Request
	from Country
}

type PerceptResponse struct {
	Request
	events []Request
}
