package gopolitical

type Event interface {
}

type MarketInteractionEvent struct {
	DateTransaction int          `json:"dateTransaction"`
	ResourceType    ResourceType `json:"resourceType"`
	Amount          float64      `json:"amount"`
	Price           float64      `json:"price"`
	Buyer           *Country     `json:"buyer"`
	Seller          *Country     `json:"seller"`
}

type CountryEvent struct {
	Event `json:"event"`
	Day   int `json:"day"`
}

type TransferResourceEvent struct {
	CountryEvent
	From     string       `json:"from"`
	To       string       `json:"to"`
	Resource ResourceType `json:"resource"`
	Amount   float64      `json:"amount"`
}

type BuyResourceEvent struct {
	CountryEvent
	From     string       `json:"from"`
	To       string       `json:"to"`
	Resource ResourceType `json:"resource"`
	Amount   float64      `json:"amount"`
	Price    float64      `json:"price"`
}

type SellResourceEvent struct {
	CountryEvent
	From     string       `json:"from"`
	To       string       `json:"to"`
	Resource ResourceType `json:"resource"`
	Amount   float64      `json:"amount"`
	Price    float64      `json:"price"`
}
