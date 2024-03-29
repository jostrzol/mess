package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/jostrzol/mess/pkg/cmd"
	"github.com/jostrzol/mess/pkg/rules"
)

func cmdError(format string, a ...any) {
	format = fmt.Sprintf("error: %s\n", format)
	fmt.Printf(format, a...)
	flag.Usage()
	os.Exit(1)
}

func runError(format string, a ...any) {
	format = fmt.Sprintf("error: %s\n", format)
	fmt.Printf(format, a...)
	os.Exit(2)
}

func main() {
	var rulesFilename = flag.String("rules", "", "path to a rules file")
	flag.Parse()

	if *rulesFilename == "" {
		cmdError("no rules file")
	}

	game, err := rules.DecodeRulesFromOs(*rulesFilename, true)
	if err != nil {
		runError("loading game rules: %s", err)
	}

	winner, err := cmd.Run(game, os.Stdin, os.Stdout)
	if errors.Is(err, cmd.ErrEOT) {
		os.Exit(3)
	} else if err != nil {
		runError("running game: %s", err)
	}

	fmt.Println()
	fmt.Println("Game over")

	if winner == nil {
		fmt.Printf("Draw!\n")
	} else {
		fmt.Printf("Winner is %v!\n", winner)
	}
}
