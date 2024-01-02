package main

import (
	"log"

	agt "github.com/BaptisteBuvron/gopolitical/server/gopolitical"
)

func main() {
	simulation, err := agt.LoadSimulation("resources/data.json")
	if err != nil {
		log.Printf("[Main] Error: %s", err)
	}
	simulation.Start()
}
