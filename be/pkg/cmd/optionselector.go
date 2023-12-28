package cmd

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/jostrzol/mess/pkg/utils"
)

func (t *interactor) selectOptions(optionTree *mess.OptionNode) (result []mess.Option, err error) {
	currentNode := optionTree
	for currentNode != nil {
		var datum mess.IOptionDatum
		datum, err = t.selectDatum(currentNode)
		if err != nil {
			return
		} else if datum == nil {
			// selected == null => move should be performed without action
			return nil, nil
		}

		currentNode, err = t.selectChild(datum.NonEmptyChildren())
		if err != nil {
			return
		}
		result = append(result, datum.IOption())
	}
	return
}

func (t *interactor) selectDatum(node *mess.OptionNode) (mess.IOptionDatum, error) {
	selector := &optionSelector{interactor: t}
	node.Accept(selector)
	if selector.err != nil {
		return nil, selector.err
	}
	return selector.result, nil
}

type optionSelector struct {
	interactor *interactor
	result     mess.IOptionDatum
	err        error
}

func (o *optionSelector) VisitPieceTypeData(message string, data mess.PieceTypeOptionData) {
	fmt.Printf("%s:\n", message)
	o.result, o.err = selectStringer(o.interactor, data.OptionData)
}

func (o *optionSelector) VisitSquareData(message string, data mess.SquareOptionData) {
	var square board.Square

	fmt.Printf("%s:\n", message)
	square, o.err = o.interactor.selectSquare()
	if o.err != nil {
		return
	}

	for _, datum := range data.OptionData {
		if square == datum.Option.Square {
			o.result = datum
			return
		}
	}
	o.err = fmt.Errorf("invalid option")
}

func (o *optionSelector) VisitMoveData(_ string, data mess.MoveOptionData) {
	o.result, o.err = o.interactor.selectMove(data)
}

func (o *optionSelector) VisitUnitData(_ string, data mess.UnitOptionData) {
	o.result = utils.Single(data.OptionData)
}

type Option = interface {
	comparable
	mess.Option
}

func (t *interactor) selectChild(children []*mess.OptionNode) (*mess.OptionNode, error) {
	switch len(children) {
	case 0:
		// result == nil => move should be performed without action
		return nil, nil
	case 1:
		return utils.Single(children), nil
	default:
		fmt.Println("Choose action:")
		var messages []string
		for _, child := range children {
			messages = append(messages, child.Message)
		}
		i, err := t.selectString(messages)
		if err != nil {
			return nil, err
		}
		return children[i], nil
	}
}

func (t *interactor) selectMove(data mess.MoveOptionData) (*mess.OptionDatum[mess.MoveOption], error) {
	println("Choose a square with your piece")
	piece, err := t.selectOwnPiece()
	if err != nil {
		return nil, err
	}

	validForPiece := make(map[board.Square]*mess.OptionDatum[mess.MoveOption], 0)
	for _, datum := range data.OptionData {
		vec := datum.Option.SquareVec
		if vec.From == piece.Square() {
			validForPiece[vec.To] = datum
		}
	}

	if len(validForPiece) == 0 {
		return nil, ErrNoMoves
	}
	println("Valid destinations:")
	for _, datum := range validForPiece {
		fmt.Printf("-> %v\n", &datum.Option.SquareVec.To)
	}

	println("Choose a destination square")
	dst, err := t.selectSquare()
	if err != nil {
		return nil, err
	}
	datum, ok := validForPiece[dst]
	if !ok {
		return nil, fmt.Errorf("invalid move")
	}

	return datum, nil
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

func selectStringer[T fmt.Stringer](t *interactor, options []T) (T, error) {
	optionStrings := make([]string, len(options))
	for i, option := range options {
		optionStrings[i] = option.String()
	}

	i, err := t.selectString(optionStrings)
	if err != nil {
		var zero T
		return zero, err
	}

	return options[i], nil
}

func (t *interactor) selectString(strings []string) (int, error) {
	for i, str := range strings {
		fmt.Printf("%v. %v\n", i+1, str)
	}
	choiceStr, err := t.scan()
	if err != nil {
		return 0, err
	}

	var i int
	_, err = fmt.Sscanf(choiceStr, "%d", &i)
	if err != nil {
		return 0, err
	}

	if i < 1 || i > len(strings) {
		return 0, fmt.Errorf("invalid option")
	}
	return i - 1, nil
}
