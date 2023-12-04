package gopolitical

type Environment struct {
	Agents map[string]AgentI
	Market Market
}

func NewEnvironment(agents map[string]AgentI, prices Prices) Environment {
	return Environment{agents, NewMarket(prices)}
}
