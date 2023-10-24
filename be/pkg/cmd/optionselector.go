package cmd

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/mess"
	"golang.org/x/exp/maps"
)

func (t *interactor) selectOptionSet(optionSets [][]mess.Option) ([]mess.Option, error) {
	if len(optionSets) == 0 {
		return []mess.Option{}, nil
	} else if len(optionSets) == 1 && len(optionSets[0]) == 0 {
		return optionSets[0], nil
	}

	for i := 0; i < len(optionSets[0]); i++ {
		groups := mess.GroupOptions(optionSets, i)

		group := maps.Values(groups)[0]
		if len(groups) > 1 {
			fmt.Println("Choose action:")
			messages := maps.Keys(groups)
			message, err := t.selectString(messages)
			if err != nil {
				return nil, err
			}
			group = groups[message]
		}

		option, err := t.selectOption(group)
		if err != nil {
			return nil, err
		}

		newOptionSets := make([][]mess.Option, 0)
		for _, options := range optionSets {
			if options[i] == option {
				newOptionSets = append(newOptionSets, options)
			}
		}
		optionSets = newOptionSets
	}

	if len(optionSets) != 1 {
		return nil, fmt.Errorf("expected exactly 1 option after selection, got: %v", len(optionSets))
	}
	return optionSets[0], nil
}

func (t *interactor) selectOption(group mess.OptionGroup) (mess.Option, error) {
	selector := &optionSelector{interactor: t}
	group.Accept(selector)
	if selector.err != nil {
		return nil, selector.err
	}
	return selector.result, nil
}

type optionSelector struct {
	interactor *interactor
	result     mess.Option
	err        error
}

func (o *optionSelector) VisitPieceTypeOptions(options []*mess.PieceTypeOption) {
	o.result, o.err = selectWithNumber(o.interactor, options)
}

func (o *optionSelector) VisitSquareOptions(options []*mess.SquareOption) {
	var square board.Square

	message := options[0].Message()
	fmt.Printf("%s:\n", message)
	square, o.err = o.interactor.selectSquare()
	if o.err != nil {
		return
	}

	for _, option := range options {
		if square == option.Square {
			o.result = option
			return
		}
	}
	o.err = fmt.Errorf("invalid option")
}

func (o *optionSelector) VisitMoveOptions(options []*mess.MoveOption) {
	var move *mess.Move
	move, o.err = o.interactor.selectMove()
	options[0].Move = move
	o.result = options[0]
}

func (o *optionSelector) VisitUnitOptions(options []*mess.UnitOption) {
	o.result = options[0]
}

func selectWithNumber[T mess.Option](t *interactor, options []T) (mess.Option, error) {
	optionsByString := make(map[string]T, 1)
	for _, option := range options {
		optionsByString[option.String()] = option
	}

	message := options[0].Message()
	fmt.Printf("%s:\n", message)

	optionStrings := maps.Keys(optionsByString)
	optionString, err := t.selectString(optionStrings)
	if err != nil {
		return nil, err
	}

	return optionsByString[optionString], nil
}

func (t *interactor) selectMove() (*mess.Move, error) {
	println("Choose a square with your piece")
	piece, err := t.selectOwnPiece()
	if err != nil {
		return nil, err
	}

	validForPiece := make(map[board.Square]mess.MoveGroup, 0)
	for _, moveGroup := range t.game.ValidMoves() {
		if moveGroup.Piece == piece {
			validForPiece[moveGroup.To] = moveGroup
		}
	}

	if len(validForPiece) == 0 {
		return nil, ErrNoMoves
	}
	println("Valid destinations:")
	for _, validMove := range validForPiece {
		fmt.Printf("-> %v\n", &validMove.To)
	}

	println("Choose a destination square")
	dst, err := t.selectSquare()
	if err != nil {
		return nil, err
	}
	moveGroup, ok := validForPiece[dst]
	if !ok {
		return nil, fmt.Errorf("invalid move")
	}

	options, err := t.selectOptionSet(moveGroup.OptionSets())
	if err != nil {
		return nil, err
	}
	move := moveGroup.Move(options)

	return move, nil
}

func (t *interactor) selectOwnPiece() (*mess.Piece, error) {
	square, err := t.selectSquare()
	if err != nil {
		return nil, err
	}

	piece, _ := t.game.Board().At(square)
	if piece == nil {
		return nil, fmt.Errorf("square is empty")
	} else if piece.Owner() != t.game.CurrentPlayer() {
		return nil, fmt.Errorf("piece belongs to the opponent")
	}

	return piece, nil
}

func (t *interactor) selectSquare() (board.Square, error) {
	squareStr, err := t.scan()
	if err != nil {
		return board.Square{}, err
	}
	square, err := board.NewSquare(squareStr)
	if err != nil {
		return board.Square{}, err
	} else if !t.game.Board().Contains(square) {
		return board.Square{}, fmt.Errorf("Square not on board")
	}
	return square, nil
}

func (t *interactor) selectString(strings []string) (string, error) {
	for i, str := range strings {
		fmt.Printf("%v. %v\n", i+1, str)
	}
	choiceStr, err := t.scan()
	if err != nil {
		return "", err
	}

	var i int
	_, err = fmt.Sscanf(choiceStr, "%d", &i)
	if err != nil {
		return "", err
	}

	if i < 1 || i > len(strings) {
		return "", fmt.Errorf("invalid option")
	}
	return strings[i-1], nil
}
