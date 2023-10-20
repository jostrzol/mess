package cmd

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/mess"
	"golang.org/x/exp/maps"
)

type Interactor struct {
	Scanner *bufio.Scanner
}

func (t *Interactor) PreTurn(state *mess.State) {
	// generate moves first so that debug logs print before the board does
	// (the moves are cached anyway, so this computation won't get wasted)
	state.ValidMoves()

	println("Board: (uppercase - white, lowercase - black)")
	println(state.PrettyString())

}

func (t *Interactor) ChooseMove(state *mess.State, validMoves []mess.GeneratedMove) (*mess.Move, error) {
	for {
		println("Choose a square with your piece")
		print("> ")
		if !t.Scanner.Scan() {
			return nil, t.Scanner.Err()
		}
		srcStr := t.Scanner.Text()
		if srcStr == "" {
			continue
		}

		src, err := board.NewSquare(srcStr)
		if err != nil {
			fmt.Printf("%v\n", err)
			println("Try again")
			continue
		} else if !state.Board().Contains(src) {
			println("Square not on board")
			println("Try again")
			continue
		}

		piece, _ := state.Board().At(src)
		if piece == nil {
			println("That square is empty!")
			continue
		} else if piece.Owner() != state.CurrentPlayer() {
			println("That belongs to your opponent!")
			continue
		}

		validForPiece := make(map[board.Square]mess.GeneratedMove, 0)
		for _, validMove := range validMoves {
			if validMove.Piece == piece {
				validForPiece[validMove.To] = validMove
			}
		}

		if len(validForPiece) == 0 {
			println("No valid moves for this piece!")
			continue
		}
		println("Valid destinations:")
		for _, validMove := range validForPiece {
			fmt.Printf("-> %v\n", &validMove.To)
		}

		println("Choose a destination square")
		print("> ")
		if !t.Scanner.Scan() {
			return nil, t.Scanner.Err()
		}
		dstStr := t.Scanner.Text()
		if dstStr == "" {
			continue
		}

		dst, err := board.NewSquare(dstStr)
		if err != nil {
			fmt.Printf("%v\n", err)
			println("Try again")
		}

		generatedMove, ok := validForPiece[dst]
		if !ok {
			println("Invalid move!")
			continue
		}

		optionSet, err := t.chooseOptionSet(generatedMove.OptionSets, 0)
		if err != nil {
			fmt.Printf("%v\n", err)
			println("Try again")
		} else if errors.Is(err, errCancel) {
			continue
		}

		move := generatedMove.ToMove(optionSet)

		return &move, nil
	}
}

func (t *Interactor) chooseOptionSet(optionSets [][]mess.Option, choiceIdx int) ([]mess.Option, error) {
	switch len(optionSets) {
	case 0:
		return nil, nil
	case 1:
		return optionSets[0], nil
	}

	optionSetsMap := make(map[string][][]mess.Option, 0)
	for _, optionSet := range optionSets {
		option := optionSet[choiceIdx]
		matchingOptionSets, found := optionSetsMap[option.String()]
		if !found {
			matchingOptionSets = [][]mess.Option{}
		}
		optionSetsMap[option.String()] = append(matchingOptionSets, optionSet)
	}

	optionStrings := maps.Keys(optionSetsMap)

	println("Choose option:")
	for i, optionString := range optionStrings {
		fmt.Printf("%v. %v\n", i+1, optionString)
	}
	print("> ")

	if !t.Scanner.Scan() {
		return nil, t.Scanner.Err()
	}
	choiceStr := t.Scanner.Text()
	if choiceStr == "" {
		return nil, errCancel
	}

	var i int
	_, err := fmt.Sscanf(choiceStr, "%d", &i)
	if err != nil {
		return nil, err
	}

	if i < 1 || i > len(optionStrings) {
		return nil, fmt.Errorf("invalid option")
	}
	optionString := optionStrings[i-1]

	narrowedOptions := optionSetsMap[optionString]

	return t.chooseOptionSet(narrowedOptions, choiceIdx+1)
}

var errCancel = fmt.Errorf("cancel")
