package gopolitical

import "time"

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
	MarketRequest
	to             *Country
	amountExecuted float64
	cost           float64
}

type MarketSellResponse struct {
	MarketRequest
	to             *Country
	amountExecuted float64
	gain           float64
}

type MarketInteraction struct {
	DateTransaction time.Time    `json:"dateTransaction"`
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
	events []Request
}
