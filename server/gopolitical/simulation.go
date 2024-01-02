package gopolitical

import (
	"log"
	"sync"
	"time"
)

type Simulation struct {
	SecondByDay float64         `json:"secondByDay"`
	Environment *Environment    `json:"environment"`
	CurrentDay  int             `json:"currentDay"`
	WebSocket   *WebSocket      `json:"-"`
	WgStart     *sync.WaitGroup `json:"-"`
	WgMiddle    *sync.WaitGroup `json:"-"`
	WgEnd       *sync.WaitGroup `json:"-"`
}

func NewSimulation(
	worldWidth int,
	worldHeight int,
	secondByDay float64,
	prices Prices,
	countries map[string]*Country,
	territories []*Territory,
	consumptionsByHabitant map[ResourceType]float64,
) Simulation {
	WgStart := new(sync.WaitGroup)
	WgEnd := new(sync.WaitGroup)
	WgMiddle := new(sync.WaitGroup)
	for _, c := range countries {
		c.WgStart = WgStart
		c.WgMiddle = WgMiddle
		c.WgEnd = WgEnd
	}
	return Simulation{secondByDay, NewEnvironment(worldWidth, worldHeight, countries, territories, prices, consumptionsByHabitant), 0, nil, WgStart, WgMiddle, WgEnd}
}

func (s *Simulation) Start() {
	// Launch all agents and added a channel to the environment
	s.WebSocket = NewWebSocket(s)
	go s.WebSocket.Start() // TODO: attendre le lancement pour commencer
	time.Sleep(1)

	log.Printf("[Simulation] Statistiques")
	log.Printf("[Simulation] Pays                     : %d", len(s.Environment.Countries))
	log.Printf("[Simulation] Territoires              : %d", len(s.Environment.Market.Env.World.Territories()))

	for _, country := range s.Environment.Countries {
		log.Printf("[Simulation] %-24s : %d", country.Name, len(country.Territories))
	}

	go s.Environment.Start()

	s.WgStart.Add(1)
	s.WgMiddle.Add(len(s.Environment.Countries))
	s.WgEnd.Add(1)

	// Start agent threads
	for _, country := range s.Environment.Countries {
		go country.Start()
	}

	for {
		// Increment current day
		s.incrementDay()
		log.Printf("[Simulation] Commencement du jour %d", s.CurrentDay)

		// Start all agents
		// TODO: Comment géré l'ajout d'un nouveau pays
		s.WgStart.Done()                             // On Commence les threads
		s.WgMiddle.Wait()                            // On attends qu'ils se termine tous
		s.WgMiddle.Add(len(s.Environment.Countries)) // On remonte le conteur de WgMiddle
		s.WgStart.Add(1)                             // On reverrouille le commencement
		s.WgEnd.Done()                               // On lance le réarmement des threads
		s.WgMiddle.Wait()                            // On s'assure qu'ils soit tous réarmé
		s.WgMiddle.Add(len(s.Environment.Countries)) // On remonte le conteur de WgMiddle
		s.WgEnd.Add(1)                               // On reverrouille la fin
		log.Println("[Simulation] Fin des actions des pays")

		// On fait correspondre les ordres d'achats et de ventes
		s.Environment.Market.HandleRequests()

		// Mettre à jour les stocks des territoires à partir des variations
		s.Environment.UpdateStocksFromVariation()

		// Mettre à jour les stocks des territoires à partir des consommations des habitants
		s.Environment.UpdateStocksFromConsumption()
		s.Environment.KillHungryHabitants()
		s.Environment.BirthHabitants()

		//Add history
		s.Environment.UpdateStockHistory(s.CurrentDay)
		s.Environment.UpdateMoneyHistory(s.CurrentDay)
		s.Environment.UpdateHabitantsHistory(s.CurrentDay)
		s.Environment.Percept = s.Environment.Market.Percept
		s.Environment.Market.Percept = make(map[string][]Request)

		//Send update to the websocket
		s.WebSocket.SendUpdate()

		log.Printf("[Simulation] Fin du jour %d", s.CurrentDay)

		// Espace dans les logs
		log.Printf("")
		log.Printf("")
		log.Printf("")

		//Wait the other day
		time.Sleep(time.Duration(s.SecondByDay) * time.Second)
		// Reader
		//var keepGoing string
		//log.Printf("[Simulation] [WAIT] Continuer ? %d", s.CurrentDay)
		//fmt.Scanln(&keepGoing)
	}
}

func (s *Simulation) incrementDay() {
	s.CurrentDay++
	s.Environment.Market.currentDay++
}
