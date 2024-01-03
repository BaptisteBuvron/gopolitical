package gopolitical

import (
	"encoding/json"
)

type Pos struct {
	X int
	Y int
}

type World struct {
	territories   map[Pos]*Territory
	maritimeAreas map[Pos]*MaritimeArea
	width         int
	height        int
}

type MaritimeArea struct {
	cover     map[Pos]bool
	neighbors map[Pos]*Territory
}

func NewMaritimeArea() *MaritimeArea {
	return &MaritimeArea{make(map[Pos]bool), make(map[Pos]*Territory)}
}

func NewWorld(territories []*Territory, width int, height int) *World {
	Debug("World", "Pré-calcul des voisins")
	// Creation de la structure
	w := World{}

	// Remplissage des champs
	w.width = width
	w.height = height
	w.territories = make(map[Pos]*Territory)
	w.maritimeAreas = make(map[Pos]*MaritimeArea)

	// On cartographie la carte pour avoir un accès en O(N)
	for _, territory := range territories {
		if territory.X < 0 || territory.X >= w.width {
			panic("Invalid territory w")
		}
		if territory.Y < 0 || territory.Y >= w.height {
			panic("Invalid territory h")
		}
		w.territories[Pos{territory.X, territory.Y}] = territory
	}

	// On pré-calcule
	unexploredMaritimeAreas := make(map[Pos]bool)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			territory := w.GetTerritoryAt(x, y)
			if territory == nil {
				unexploredMaritimeAreas[Pos{x, y}] = true
			}
		}
	}
	// On trouve toute les zones
	for len(unexploredMaritimeAreas) != 0 {
		// Création d'une nouvelle zone
		area := NewMaritimeArea()

		// On prend la première position
		var current Pos
		for current = range unexploredMaritimeAreas {
			break
		}

		// On crée un pile pour stocker temporairement les positions
		queue := make([]Pos, 1)
		queue[0] = current
		delete(unexploredMaritimeAreas, current)

		// On itère jusqu'a qu'il n'y en ai plus
		for len(queue) != 0 {
			// On supprime le premier dans la pile
			current = queue[0]
			queue = queue[1:]

			// On ajoute cette zone maritime
			w.maritimeAreas[current] = area
			area.cover[current] = true

			// On marque ces voisins pour être traité
			for _, near := range w.Near(current) {
				if unexploredMaritimeAreas[near] == true {
					delete(unexploredMaritimeAreas, near)
					queue = append(queue, near)
				}
			}
		}

		// Find all neighbor in area
		for pos := range area.cover {
			for _, near := range w.Near(pos) {
				neighbor := w.territories[near]
				if neighbor != nil {
					area.neighbors[near] = neighbor
				}
			}
		}
	}
	// On vérifie que tout les points sont valides
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if w.GetTerritoryAt(x, y) == nil && w.GetMaritimeAreaAt(x, y) == nil {
				panic("Void detected")
			}
			if w.GetTerritoryAt(x, y) != nil && w.GetMaritimeAreaAt(x, y) != nil {
				panic("Quantum superposition detected")
			}
		}
	}

	// Return created world
	Debug("World", "Pré-calcul terminé")
	return &w
}

func (w *World) Delta(pos Pos, x int, y int) Pos {
	return Pos{PositiveModulus(pos.X+x, w.width), PositiveModulus(pos.Y+y, w.height)}
}

func (w *World) Territories() map[Pos]*Territory {
	return w.territories
}

type PartialWorld struct {
	Width       int          `json:"width"`
	Height      int          `json:"height"`
	Territories []*Territory `json:"territories"`
}

func (w *World) MarshalJSON() ([]byte, error) {
	// Insert the string directly into the Data member
	territories := make([]*Territory, 0)
	for _, territory := range w.Territories() {
		territories = append(territories, territory)
	}
	a := PartialWorld{w.width, w.height, territories}
	data, err := json.Marshal(a)
	return data, err
}

func (w *World) Near(pos Pos) []Pos {
	return []Pos{
		w.Delta(pos, -1, 0),
		w.Delta(pos, 0, -1),
		w.Delta(pos, 1, 0),
		w.Delta(pos, 0, 1),
	}
}

func (w *World) GetTerritoryAt(x int, y int) *Territory {
	return w.territories[Pos{PositiveModulus(x, w.width), PositiveModulus(y, w.height)}]
}

func (w *World) GetMaritimeAreaAt(x int, y int) *MaritimeArea {
	return w.maritimeAreas[Pos{PositiveModulus(x, w.width), PositiveModulus(y, w.height)}]
}

func (w *World) FindNeighborTerritoriesOfCountry(country *Country) map[Pos]*Territory {
	// Trouver tous les territoires d'un pays
	territories := make(map[Pos]*Territory)

	// On itère sur tous ses territoire
	for _, countryTerritory := range country.Territories {

		// On continue le traitement sur ces voisins
		pos := Pos{countryTerritory.X, countryTerritory.Y}
		for _, near := range w.Near(pos) {
			neighborTerritory := w.territories[near]
			if neighborTerritory == nil { // Si il n'y a pas de territoire, c'est une ouverture sur la mer
				// On marque tous les voisins à cette zones
				area := w.maritimeAreas[near]
				for posNeighbor, neighbor := range area.neighbors {
					if neighbor.Country.ID != country.ID {
						territories[posNeighbor] = neighbor
					}
				}
				// Si il un territoire et qu'il n'appartient pas au pays en lui même
			} else if neighborTerritory.Country.ID != country.ID {
				// On le marque en tant que voisins du pays
				territories[near] = neighborTerritory
			}
		}
	}
	return territories
}

func (w *World) FindNeighborTerritoriesOfCountryWith(country *Country, other string) map[Pos]*Territory {
	territories := make(map[Pos]*Territory)
	for pos, territory := range w.FindNeighborTerritoriesOfCountry(country) {
		if territory.Country.ID == other {
			territories[pos] = territory
		}
	}
	return territories
}
