package gopolitical

type Territory struct {
	X                int                              `json:"x"`
	Y                int                              `json:"y"`
	Name             string                           `json:"name"`
	Variations       []Variation                      `json:"variations"`
	Stock            map[ResourceType]float64         `json:"stock"`
	StockHistory     map[int]map[ResourceType]float64 `json:"stockHistory"`
	Country          *Country                         `json:"country"`
	Habitants        int                              `json:"habitants"`
	HabitantsHistory map[int]int                      `json:"habitantsHistory"`
}

func NewTerritory(x int, y int, name string, variations []Variation, stock map[ResourceType]float64, country *Country, habitant int) *Territory {
	return &Territory{x, y, name, variations, stock, make(map[int]map[ResourceType]float64), country, habitant, make(map[int]int)}
}

func (t Territory) Start() {

}

func (t Territory) Percept() {

}

func (t Territory) Deliberate() {

}

func (t Territory) Act() {

}

func (t *Territory) GetSurplus(daysToSecure float64) map[ResourceType]float64 {
	surplus := make(map[ResourceType]float64)
	//On garde 3 jours de surplus

	for resource, consumption := range t.Country.consumptionByHabitant {
		surplusAmount := t.Stock[resource] - (float64(t.Habitants)*consumption)*daysToSecure
		if surplusAmount > 0 {
			surplus[resource] = surplusAmount
		}
	}
	return surplus
}
