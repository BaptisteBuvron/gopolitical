package gopolitical

type ResourceType string

type Resource struct {
	Id       int          `json:"id"`
	Name     ResourceType `json:"name"`
	Quantity int          `json:"quantity"`
}

const (
	ARMAMENT ResourceType = "armement"
)
const (
	ARMAMENT_NEEDED_BY_TERRITORY = 100
	DAYS_TO_SECURE               = 5
	DAYS_TO_WARS                 = 3
)
