package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jostrzol/mess/config"
)

func main() {
	var configFilename = flag.String("rules", "./rules.hcl", "path to a rules config file")
	flag.Parse()

	config, err := config.DecodeConfig(*configFilename)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	log.Printf("Configuration is %#v", config)

	state, err := config.ToGame()
	if err != nil {
		log.Fatalf("Failed to create initial game state: %s", err)
	}
	log.Printf("Initial game state loaded")

	controller := config.ToController()

	winner, err := controller.DecideWinner(state)
	if err != nil {
		log.Fatalf("Failed to resolve game: %s", err)
	}

	if winner == nil {
		fmt.Printf("Draw!")
	} else {
		fmt.Printf("Winner is %v!", winner)
	}
}
