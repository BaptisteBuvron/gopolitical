package gopolitical

type Agent struct {
	ID   string
	Name string
}

type AgentI interface {
	Start()
	Percept()
	Deliberate()
	Act()
}
