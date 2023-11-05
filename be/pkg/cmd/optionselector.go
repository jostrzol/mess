package cmd

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/jostrzol/mess/pkg/utils"
	"golang.org/x/exp/maps"
)

func (t *interactor) selectOptions(optionTree mess.OptionTree) (result []mess.Option, err error) {
	selected := &selected{node: optionTree}
	for selected.node != nil {
		selected, err = t.selectOption(optionTree)
		if err != nil {
			return
		} else if selected == nil {
			// selected == null => move should be performed without action
			return nil, nil
		}
		result = append(result, selected.option)
	}
	return
}

func (t *interactor) selectOption(optionTree mess.OptionTree) (*selected, error) {
	selector := &optionSelector{interactor: t}
	optionTree.Accept(selector)
	if selector.err != nil {
		return nil, selector.err
	}
	return selector.result, nil
}

type optionSelector struct {
	interactor *interactor
	result     *selected
	err        error
}

type selected struct {
	option mess.Option
	node   mess.OptionTree
}

func (o *optionSelector) VisitPieceTypeNode(options map[*mess.PieceTypeOption]mess.OptionTree) {
	o.result, o.err = selectWithNumber(o.interactor, options)
}

func (o *optionSelector) VisitSquareNode(options map[*mess.SquareOption]mess.OptionTree) {
	var square board.Square

	message := maps.Keys(options)[0].Message()
	fmt.Printf("%s:\n", message)
	square, o.err = o.interactor.selectSquare()
	if o.err != nil {
		return
	}

	for option, node := range options {
		if square == option.Square {
			o.result = &selected{option, node}
			return
		}
	}
	o.err = fmt.Errorf("invalid option")
}

func (o *optionSelector) VisitMoveNode(options map[*mess.MoveOption]mess.OptionTree) {
	var move *mess.Move
	move, o.err = o.interactor.selectMove()
	if o.err != nil {
		return
	}
	option, node := utils.SingleEntry(options)
	option.Move = move
	o.result = &selected{option, node}
}

func (o *optionSelector) VisitUnitNode(options map[*mess.UnitOption]mess.OptionTree) {
	option, node := utils.SingleEntry(options)
	o.result = &selected{option, node}
}

func (o *optionSelector) VisitMessageNode(options map[string]mess.OptionTree) {
	switch len(options) {
	case 0:
		// o.result == null => move should be performed without action
	case 1:
		_, node := utils.SingleEntry(options)
		node.Accept(o)
	default:
		fmt.Println("Choose action:")
		messages, nodes := utils.Entries(options)
		i, err := o.interactor.selectString(messages)
		if err != nil {
			o.err = err
			return
		}
		nodes[i].Accept(o)
	}
}

type Option = interface {
	comparable
	mess.Option
}

func selectWithNumber[T Option](t *interactor, options map[T]mess.OptionTree) (*selected, error) {
	optionsByString := make(map[string]*selected, 1)
	var message string
	for option, node := range options {
		optionsByString[option.String()] = &selected{option, node}
		message = option.Message()
	}

	fmt.Printf("%s:\n", message)

	optionStrings := maps.Keys(optionsByString)
	i, err := t.selectString(optionStrings)
	if err != nil {
		return nil, err
	}

	return optionsByString[optionStrings[i]], nil
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

	options, err := t.selectOptions(moveGroup.OptionTree())
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
