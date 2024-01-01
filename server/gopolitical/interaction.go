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
	BuyID      int
	from       *Country
	territoire *Territory
	resources  ResourceType
	amount     float64
}

type MarketSellRequest struct {
	MarketRequest
	SellID     int
	from       *Country
	territoire *Territory
	resources  ResourceType
	amount     float64
}

type MarketBuyResponse struct {
	Request
	Event          `json:"event"`
	Day            int          `json:"day"`
	ResourceType   ResourceType `json:"resourceType"`
	From           string       `json:"from"`
	AmountExecuted float64      `json:"amountExecuted"`
	Cost           float64      `json:"cost"`
}

type MarketSellResponse struct {
	Request
	Event          `json:"event"`
	Day            int          `json:"day"`
	ResourceType   ResourceType `json:"resourceType"`
	To             string       `json:"to"`
	AmountExecuted float64      `json:"amountExecuted"`
	Gain           float64      `json:"gain"`
}

type MarketInteraction struct {
	DateTransaction int          `json:"dateTransaction"`
	ResourceType    ResourceType `json:"resourceType"`
	Amount          float64      `json:"amount"`
	Price           float64      `json:"price"`
	Buyer           *Country     `json:"buyer"`
	Seller          *Country     `json:"seller"`
}

//Percept

type PerceptRequest struct {
	Request
	from *Country
}

type PerceptResponse struct {
	Request
	events          []Request
	RelationManager *RelationManager
	World           *World
}

type AttackRequest struct {
	Request
	to       *Territory
	armement float64
}

type AttackResponse struct {
	Request
	to *Country
}
