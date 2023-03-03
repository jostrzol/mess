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

	state, controller, err := config.DecodeConfig(*configFilename)
	if err != nil {
		log.Fatalf("loading game rules: %s", err)
	}

	print(state.String())

	winner := controller.DecideWinner(state)

	if winner == nil {
		fmt.Printf("Draw!")
	} else {
		fmt.Printf("Winner is %v!", winner)
	}
}
