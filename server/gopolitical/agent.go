package gopolitical

type Agent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AgentI interface {
	Start()
	Percept()
	Deliberate()
	Act()
}
