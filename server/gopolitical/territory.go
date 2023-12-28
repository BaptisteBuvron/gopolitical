package gopolitical

type Territory struct {
	X          int                      `json:"x"`
	Y          int                      `json:"y"`
	Variations []Variation              `json:"variations"`
	Stock      map[ResourceType]float64 `json:"stock"`
	Country    *Country                 `json:"country"`
	Habitants  int                      `json:"habitants"`
}

func NewTerritory(x int, y int, variations []Variation, stock map[ResourceType]float64, country *Country, habitant int) *Territory {
	return &Territory{x, y, variations, stock, country, habitant}
}

func (t Territory) Start() {

}

func (t Territory) Percept() {

}

func (t Territory) Deliberate() {

}

func (t Territory) Act() {

}

func (t *Territory) GetSurplus() map[ResourceType]float64 {
	//TODO: Rendre générique pour toutes les ressources
	surplus := make(map[ResourceType]float64)
	//On garde 3 jours de surplus
	surplusFood := t.Stock["food"] - (float64(t.Habitants)*FOOD_BY_HABITANT)*3
	surplusWater := t.Stock["water"] - (float64(t.Habitants)*WATER_BY_HABITANT)*3

	if surplusFood > 0 {
		surplus["food"] = surplusFood
	}
	if surplusWater > 0 {
		surplus["water"] = surplusWater
	}
	return surplus
}
