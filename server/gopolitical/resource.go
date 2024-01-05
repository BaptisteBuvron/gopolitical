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
	ARMAMENT_NEEDED_BY_TERRITORY         = 100
	ARMAMENT_NEEDED_FOR_WAR_BY_TERRITORY = 300
	DAYS_TO_SECURE                       = 4
	DAYS_TO_WARS                         = 2
	BIRTH_RATIO                          = 0.03
	STARVATION_RATIO                     = 0.1
	INCOMPRESSIBLE_CAPTURE_RATIO         = 0.2
	RELATION_RATIO_DEFEND                = 2 / 3
	RELATION_RATIO_ATTACK                = 1 / 3
)
