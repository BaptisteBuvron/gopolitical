package main

import (
	gopolitical "github.com/BaptisteBuvron/gopolitical/server/gopolitical"
)

func main() {
	simulation, err := gopolitical.LoadSimulation("resources/data.json")
	if err != nil {
		gopolitical.Info("Main", "Error: %s", err)
	}
	simulation.Start()
}
