package gopolitical

type ResourceType string

type Ressource struct {
	Id       int          `json:"id"`
	Name     ResourceType `json:"name"`
	Quantity int          `json:"quantity"`
}
