package gopolitical

import (
	"fmt"
	"time"
)

type Simulation struct {
	SecondByDay float64      `json:"secondByDay"`
	Environment *Environment `json:"environment"`
	WebSocket   *WebSocket   `json:"-"`
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
	return Simulation{
		secondByDay,
		NewEnvironment(worldWidth, worldHeight, countries, territories, prices, consumptionsByHabitant),
		nil,
	}
}

func (s *Simulation) Start() {
	// Launch all agents and added a channel to the environment
	s.WebSocket = NewWebSocket(s)
	go s.WebSocket.Start() // TODO: attendre le lancement pour commencer
	time.Sleep(1)

	Debug("Simulation", "Statistiques")
	Debug("Simulation", "Pays                     : %d", len(s.Environment.Countries))
	Debug("Simulation", "Territoires              : %d", len(s.Environment.Market.Env.World.Territories()))

	// Start agent threads
	for _, agent := range s.Environment.Agents {
		go agent.Start()
	}

	for {
		startTimeDay := time.Now()
		Info("Simulation", "Commencement du jour %d", s.Environment.CurrentDay)

		// On demande au agents de percevoir leurs environnement de façons synchrone (C'est à eux même de s'envoyer leurs perceptions)
		for _, agent := range s.Environment.Agents {
			agent.Percept(s.Environment)
		}

		// Les agents se sont envoyer leurs perception grâce à percept (Nous avons que Country ici, mais cela permet l'extension d'autres futures agents)

		// On récupère les actions de chacun
		actions := make([]Action, 0)
		for _, agent := range s.Environment.Agents {
			agentActions := agent.Act()
			for _, action := range agentActions {
				actions = append(actions, action)
			}
		}

		Debug("Simulation", "Fin des actions des pays")

		// L'environnement traites toutes les requêtes
		s.Environment.HandleActions(actions)

		// L'environnement ce mets à jour
		s.Environment.Update()

		// Envoie de l'update dans le websocket
		s.WebSocket.SendUpdate()

		Debug("Simulation", "Fin du jour %d", s.Environment.CurrentDay)

		// Espace dans les logs
		Debug("Simulation", "")
		Debug("Simulation", "")
		Debug("Simulation", "")

		//Wait the other day
		endTimeDay := time.Now()
		expectedEndTimeDay := startTimeDay.Add(time.Duration(s.SecondByDay*1000000000) * time.Nanosecond)
		if expectedEndTimeDay.After(endTimeDay) {
			time.Sleep(expectedEndTimeDay.Sub(endTimeDay))
		}
		// Attente optionnelle
		if DebugEnabled() {
			var keepGoing string
			Debug("Simulation", "Entrer pour continuer")
			fmt.Scanln(&keepGoing)
		}
	}
}
