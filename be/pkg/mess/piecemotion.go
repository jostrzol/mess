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
	Generate         MoveGeneratorFunc
	ChoiceGenerators []ChoiceFunc
	Action           MoveActionFunc
}

type MoveGeneratorFunc = func(*Piece) []board.Square
type ChoiceFunc = func(*Piece, board.Square, board.Square) Choice
type MoveActionFunc = func(*Piece, board.Square, board.Square, []Option) error

type chainMotions []Motion

func (g chainMotions) Generate(piece *Piece) []GeneratedMove {
	resultMap := make(map[brd.Square]GeneratedMove, 0)
	for _, motion := range g {
		name := motion.Name
		destinations := motion.Generate(piece)

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
				optionSets = make([][]Option, 1)
				optionSets[0] = make([]Option, 0)
			}

			resultMap[destination] = GeneratedMove{
				Name:       name,
				Piece:      piece,
				From:       source,
				To:         destination,
				Action:     motion.Action,
				OptionSets: optionSets,
			}
		}
	}

	return maps.Values(resultMap)
}

func choicesToOptionSets(choices []Choice) [][]Option {
	options := make([][]Option, len(choices))
	for i, choice := range choices {
		if choice == nil {
			return [][]Option{}
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

type GeneratedMove struct {
	Name       string
	Piece      *Piece
	From       brd.Square
	To         brd.Square
	Action     MoveActionFunc
	OptionSets [][]Option
}

func (m *GeneratedMove) FilterOptionSets(predicate func([]Option) bool) {
	result := make([][]Option, 0, len(m.OptionSets))
	for _, optionSet := range m.OptionSets {
		if predicate(optionSet) {
			result = append(result, optionSet)
		}
	}
	m.OptionSets = result
}

func (m *GeneratedMove) ToMove(optionSet []Option) Move {
	return Move{
		Name:      m.Name,
		Piece:     m.Piece,
		From:      m.From,
		To:        m.To,
		Action:    m.Action,
		OptionSet: optionSet,
	}
}

type Move struct {
	Name      string
	Piece     *Piece
	From      brd.Square
	To        brd.Square
	Action    MoveActionFunc
	OptionSet []Option
}

func (m *Move) Perform() error {
	err := m.Piece.MoveTo(m.To)
	if err != nil {
		return err
	}
	if m.Action != nil && m.OptionSet != nil {
		err = m.Action(m.Piece, m.From, m.To, m.OptionSet)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Move) String() string {
	return fmt.Sprintf("%v: %v->%v", m.Piece, m.From, m.To)
}
