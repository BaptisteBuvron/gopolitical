package main

import (
	"fmt"

	agt "github.com/BaptisteBuvron/gopolitical/server/gopolitical"
)

func main() {
	simulation, err := agt.LoadSimulation("resources/data.json")
	simulation.Start()

	if err != nil {
		fmt.Printf("Error: %s", err)
	} else {
		fmt.Printf("Simulation loaded%v\n", simulation)
	}
}
