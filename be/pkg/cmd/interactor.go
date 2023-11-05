package cmd

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/jostrzol/mess/pkg/mess"
)

const barLength = 80

var ErrNoMoves = fmt.Errorf("no valid moves for this piece")
var ErrCancel = fmt.Errorf("cancel")
var ErrEOT = fmt.Errorf("EOT")

type interactor struct {
	scanner *bufio.Scanner
	game    *mess.Game
}

func newInteractor(game *mess.Game, scanner *bufio.Scanner) *interactor {
	return &interactor{
		scanner: scanner,
		game:    game,
	}
}

func (t *interactor) Run() (*mess.Player, error) {
	var winner *mess.Player
	isFinished := false

	t.printState()
	for !isFinished {
		optionTree, err := t.game.TurnOptions()
		if err != nil {
			return nil, err
		}

		options, err := t.selectOptions(optionTree)
		if errors.Is(err, ErrCancel) {
			t.printMessage("Move cancelled")
			t.printState()
			continue
		} else if errors.Is(err, ErrEOT) {
			t.printMessage("Encountered end of text")
			return nil, err
		} else if err != nil {
			t.printMessage("-> Error: %v!", err)
			t.printState()
			continue
		}

		err = t.game.Turn(options)
		if err != nil {
			return nil, fmt.Errorf("executing turn action: %v", err)
		}

		t.game.EndTurn()
		t.printState()

		isFinished, winner = t.game.PickWinner()
	}
	return winner, nil
}

func (t *interactor) printState() {
	fmt.Println(t.game.PrettyString())
}

func (t *interactor) printMessage(format string, a ...any) {
	t.printBar()
	fmt.Print("| ")
	fmt.Printf(format, a...)
	fmt.Println()
	t.printBar()
}

func (t *interactor) printBar() {
	for i := 0; i < barLength; i++ {
		fmt.Print("-")
	}
	fmt.Println()
}

func (t *interactor) scan() (string, error) {
	fmt.Print("> ")
	if !t.scanner.Scan() {
		if t.scanner.Err() == nil {
			return "", ErrEOT
		}
		return "", t.scanner.Err()
	}
	text := t.scanner.Text()
	if text == "" {
		return "", ErrCancel
	}
	return text, nil
}
