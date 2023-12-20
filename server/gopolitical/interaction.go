package gopolitical

type Channel = chan Request

type Request interface {
}

type MarketBuyRequest struct {
	Request
	from      Country
	resources ResourceType
	amount    int
}

type MarketSellRequest struct {
	Request
	from      Country
	resources ResourceType
	amount    int
}

type PerceptRequest struct {
	Request
	from Country
}

type PerceptResponse struct {
	Request
	events []Request
}
