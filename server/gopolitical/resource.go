package gopolitical

type ResourceType string

type Resource struct {
	Id       int          `json:"id"`
	Name     ResourceType `json:"name"`
	Quantity int          `json:"quantity"`
}
