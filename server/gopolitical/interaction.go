package gopolitical

type Channel = chan any

type MarketBuyRequest struct {
	from      Country
	resources ResourceType
	amount    int
}

type MarketSellRequest struct {
	from      Country
	resources ResourceType
	amount    int
}
