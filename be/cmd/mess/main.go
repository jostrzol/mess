package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/jostrzol/mess/pkg/cmd"
	"github.com/jostrzol/mess/pkg/mess"
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

	scanner := bufio.NewScanner(os.Stdin)
	interactor := &cmd.Interactor{Scanner: scanner}
	game, err := rules.DecodeRules(*rulesFilename, true)
	if err != nil {
		runError("loading game rules: %s", err)
	}

	var winner *mess.Player
	isFinished := false
	for !isFinished {
		interactor.PreTurn(game.State)

		moves := game.ValidMoves()
		move, err := interactor.ChooseMove(game.State, moves)
		if err != nil {
			runError("choosing move: %s", err)
		}

		err = move.Perform()
		if err != nil {
			runError("performing move: %s", err)
		}

		game.EndTurn()
		isFinished, winner = game.PickWinner()
	}

	if winner == nil {
		fmt.Printf("Draw!\n")
	} else {
		fmt.Printf("Winner is %v!\n", winner)
	}
}
