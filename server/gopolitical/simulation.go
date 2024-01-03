package gopolitical

import (
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

	Debug("Simulation", "Statistiques")
	Debug("Simulation", "Pays                     : %d", len(s.Environment.Countries))
	Debug("Simulation", "Territoires              : %d", len(s.Environment.Market.Env.World.Territories()))

	for _, country := range s.Environment.Countries {
		Debug("Simulation", "%-24s : %d", country.Name, len(country.Territories))
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
		startTimeDay := time.Now()
		// Increment current day
		s.incrementDay()
		Info("Simulation", "Commencement du jour %d", s.CurrentDay)

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
		Debug("Simulation", "Fin des actions des pays")

		// On fait correspondre les ordres d'achats et de ventes
		s.Environment.Market.HandleRequests()

		// Mettre à jour les stocks des territoires à partir des variations
		s.Environment.UpdateStocksFromVariation()

		// Mettre à jour les stocks des territoires à partir des consommations des habitants
		s.Environment.UpdateStocksFromConsumption()
		s.Environment.ApplyRulesOfLife()

		//Add history
		s.Environment.UpdateStockHistory(s.CurrentDay)
		s.Environment.UpdateMoneyHistory(s.CurrentDay)
		s.Environment.UpdateHabitantsHistory(s.CurrentDay)
		s.Environment.Percept = s.Environment.Market.Percept
		s.Environment.Market.Percept = make(map[string][]Request)

		//Send update to the websocket
		s.WebSocket.SendUpdate()

		Debug("Simulation", "Fin du jour %d", s.CurrentDay)

		// Espace dans les logs
		Debug("Simulation", "")
		Debug("Simulation", "")
		Debug("Simulation", "")

		//Wait the other day
		endTimeDay := time.Now()
		expectedEndTimeDay := startTimeDay.Add(time.Duration(s.SecondByDay) * time.Second)
		if expectedEndTimeDay.After(endTimeDay) {
			time.Sleep(expectedEndTimeDay.Sub(endTimeDay))
		}
		/*
			// Attente optionnelle
			if DebugEnabled() {
				var keepGoing string
				Debug("Simulation", "Continuer ?")
				fmt.Scanln(&keepGoing)
			}
		*/
	}
}

func (s *Simulation) incrementDay() {
	s.CurrentDay++
	s.Environment.Market.currentDay++
}
