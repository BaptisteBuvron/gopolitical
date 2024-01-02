package gopolitical

import (
	"fmt"
)

// Les relations :
// Un pays choisira sa meilleurs relation pour commercer

type RelationKey struct {
	C1 string
	C2 string
}

type Relations map[string]map[string]float64

type RelationManager struct {
	relations        Relations
	relationsInverse Relations
}

func NewRelationManager() *RelationManager {
	return &RelationManager{
		make(map[string]map[string]float64),
		make(map[string]map[string]float64),
	}
}

// O(1)
func (rm *RelationManager) GetRelation(country string, other string) float64 {
	if country < other {
		other, country = country, other
	}
	if relations, ok := rm.relations[country]; ok {
		if relation, ok := relations[other]; ok {
			return relation
		}
	}
	return 0.0
}

// O(1)
func (rm *RelationManager) UpdateRelation(country string, other string, value float64) {
	if country < other {
		other, country = country, other
	}
	if relations, ok := rm.relations[country]; !ok {
		relations = make(map[string]float64, 1)
		relations[other] = value
		rm.relations[country] = relations
	} else {
		relations[other] = value
	}
	if relations, ok := rm.relationsInverse[other]; !ok {
		relations = make(map[string]float64, 1)
		relations[country] = value
		rm.relationsInverse[other] = relations
	} else {
		relations[country] = value
	}
}

// O(Countries)
func (rm *RelationManager) RemoveCountry(country string) {
	delete(rm.relations, country)
	for _, relations := range rm.relations {
		delete(relations, country)
	}
	delete(rm.relationsInverse, country)
	for _, relations := range rm.relations {
		delete(relations, country)
	}
}

// O(Relations)
func (rm *RelationManager) ToString() string {
	var display string = "RelationManager["
	for country, relations := range rm.relations {
		for other, relation := range relations {
			display += country
			display += "-"
			display += other
			display += ":"
			display += fmt.Sprintf("%f ", relation)
		}
	}
	display += "]"
	return display
}
