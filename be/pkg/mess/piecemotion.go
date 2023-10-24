package mess

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/board"
	brd "github.com/jostrzol/mess/pkg/board"
	"golang.org/x/exp/maps"
)

type Motion struct {
	Name             string
	MoveGenerator    MoveGeneratorFunc
	ChoiceGenerators []MoveChoiceGeneratorFunc
	Action           MoveActionFunc
}

type MoveGeneratorFunc = func(*Piece) []board.Square
type MoveChoiceGeneratorFunc = func(*Piece, board.Square, board.Square, []Option) Choice
type MoveActionFunc = func(*Piece, board.Square, board.Square, []Option) error

type chainMotions []Motion

func (g chainMotions) Generate(piece *Piece) []MoveGroup {
	resultMap := make(map[brd.Square]MoveGroup, 0)
	for _, motion := range g {
		name := motion.Name
		destinations := motion.MoveGenerator(piece)

		for _, destination := range destinations {
			source := piece.Square()

			var choiceGenerators []func([]Option) Choice
			for _, generator := range motion.ChoiceGenerators {
				generatorCopy := generator
				generatorClosure := func(options []Option) Choice {
					return generatorCopy(piece, source, destination, options)
				}
				choiceGenerators = append(choiceGenerators, generatorClosure)
			}

			optionSets := choiceGeneratorsToOptionSets(choiceGenerators)
			resultMap[destination] = MoveGroup{
				Name:       name,
				Piece:      piece,
				From:       source,
				To:         destination,
				action:     motion.Action,
				optionSets: optionSets,
			}
		}
	}

	return maps.Values(resultMap)
}

type MoveGroup struct {
	Name       string
	Piece      *Piece
	From       brd.Square
	To         brd.Square
	action     MoveActionFunc
	optionSets [][]Option
}

func (mg MoveGroup) Length() int {
	return len(mg.optionSets)
}

func (mg MoveGroup) OptionSets() [][]Option {
	return mg.optionSets
}

func (mg MoveGroup) Moves() []*Move {
	moves := make([]*Move, 0, len(mg.optionSets))
	for _, options := range mg.optionSets {
		moves = append(moves, mg.Move(options))
	}
	return moves
}

func (mg MoveGroup) Single() *Move {
	if mg.Length() != 1 {
		err := fmt.Errorf("expected move group length of 1, got: %v", mg.Length())
		panic(err)
	}
	return mg.Move(mg.optionSets[0])
}

func (mg MoveGroup) FilterMoves(predicate func(*Move) bool) MoveGroup {
	validOptionSets := make([][]Option, 0, len(mg.optionSets))
	for _, options := range mg.optionSets {
		move := mg.Move(options)
		if predicate(move) {
			validOptionSets = append(validOptionSets, options)
		}
	}
	return MoveGroup{
		Name:       mg.Name,
		Piece:      mg.Piece,
		From:       mg.From,
		To:         mg.To,
		action:     mg.action,
		optionSets: validOptionSets,
	}
}

func (mg MoveGroup) Move(options []Option) *Move {
	return &Move{
		Name:    mg.Name,
		Piece:   mg.Piece,
		From:    mg.From,
		To:      mg.To,
		Options: options,
		action:  mg.action,
	}
}

type Move struct {
	Name    string
	Piece   *Piece
	From    brd.Square
	To      brd.Square
	Options []Option
	action  MoveActionFunc
}

func (m *Move) Perform() error {
	err := m.Piece.MoveTo(m.To)
	if err != nil {
		return err
	}
	if m.action != nil && m.Options != nil {
		err = m.action(m.Piece, m.From, m.To, m.Options)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Move) String() string {
	return fmt.Sprintf("%v(%v): %v->%v", m.Name, m.Piece, m.From, m.To)
}
