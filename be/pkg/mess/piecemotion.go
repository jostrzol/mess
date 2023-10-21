package mess

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/board"
	brd "github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/utils"
	"golang.org/x/exp/maps"
)

type Motion struct {
	Name             string
	MoveGenerator    MoveGeneratorFunc
	ChoiceGenerators []ChoiceGeneratorFunc
	Action           MoveActionFunc
}

type MoveGeneratorFunc = func(*Piece) []board.Square
type ChoiceGeneratorFunc = func(*Piece, board.Square, board.Square) Choice
type MoveActionFunc = func(*Piece, board.Square, board.Square, []Option) error

type chainMotions []Motion

func (g chainMotions) Generate(piece *Piece) []MoveGroup {
	resultMap := make(map[brd.Square]MoveGroup, 0)
	for _, motion := range g {
		name := motion.Name
		destinations := motion.MoveGenerator(piece)

		for _, destination := range destinations {
			source := piece.Square()

			var optionSets [][]Option
			if len(motion.ChoiceGenerators) > 0 {
				choices := make([]Choice, 0, len(motion.ChoiceGenerators))
				for _, choiceGenerator := range motion.ChoiceGenerators {
					choice := choiceGenerator(piece, source, destination)
					choices = append(choices, choice)
				}

				optionSets = choicesToOptionSets(choices)
			} else {
				// One empty option set (no choices, action will be performed)
				optionSets = [][]Option{[]Option{}}
			}

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

func choicesToOptionSets(choices []Choice) [][]Option {
	options := make([][]Option, len(choices))
	for i, choice := range choices {
		if choice == nil {
			// One nil option set (action won't be performed)
			return [][]Option{nil}
		}
		options[i] = choice.GenerateOptions()
	}
	result := make([][]Option, 0)
	for mi := utils.NewMultiindexLike(options); !mi.IsEnd(); mi.Next() {
		optionSet := make([]Option, 0, len(choices))
		for j, k := range mi.Current() {
			optionSet = append(optionSet, options[j][k])
		}
		result = append(result, optionSet)
	}
	return result
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

func (mg MoveGroup) ChoicesNumber() int {
	if mg.Length() == 0 {
		return 0
	}
	return len(mg.optionSets[0])
}

func (mg MoveGroup) Moves() []Move {
	moves := make([]Move, 0, len(mg.optionSets))
	for _, options := range mg.optionSets {
		moves = append(moves, mg.move(options))
	}
	return moves
}

func (mg MoveGroup) Single() Move {
	if mg.Length() != 1 {
		err := fmt.Errorf("expected move group length of 1, got: %v", mg.Length())
		panic(err)
	}
	return mg.move(mg.optionSets[0])
}

func (mg MoveGroup) FilterMoves(predicate func(Move) bool) MoveGroup {
	validOptionSets := make([][]Option, 0, len(mg.optionSets))
	for _, options := range mg.optionSets {
		move := mg.move(options)
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

func (mg MoveGroup) move(options []Option) Move {
	return Move{
		Name:    mg.Name,
		Piece:   mg.Piece,
		From:    mg.From,
		To:      mg.To,
		Options: options,
		action:  mg.action,
	}
}

func (mg MoveGroup) UniqueOptionStrings(index int) []string {
	optionStrings := make(map[string]struct{}, 1)
	for _, options := range mg.optionSets {
		option := options[index]
		optionStrings[option.String()] = struct{}{}
	}
	return maps.Keys(optionStrings)
}

func (mg MoveGroup) FilterMovesByOptionString(index int, str string) MoveGroup {
	return mg.FilterMoves(func(m Move) bool {
		return m.Options[index].String() == str
	})
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
