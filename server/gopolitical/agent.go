package gopolitical

type Action interface {
	Execute(env *Environment)
}

type Agent interface {
	Start()
	Percept(*Environment)
	Deliberate()
	Act() []Action
	GetID() string
	CleanUp(*Environment)
}
