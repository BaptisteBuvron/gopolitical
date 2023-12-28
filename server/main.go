package main

import (
	agt "github.com/BaptisteBuvron/gopolitical/server/gopolitical"
	"log"
)

func main() {
	simulation, err := agt.LoadSimulation("resources/data.json")
	if err != nil {
		log.Printf("Error: %s", err)
	} else {
		log.Printf("Simulation loaded\n")
	}
	simulation.Start()

}
