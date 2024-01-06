package mess

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/board"
	brd "github.com/jostrzol/mess/pkg/board"
	"golang.org/x/exp/maps"
)

type Motion struct {
	Name          string
	MoveGenerator MoveGeneratorFunc
	ChoiceFunc    MoveChoiceFunc
	Action        MoveActionFunc
}

type MoveGeneratorFunc = func(*Piece) []board.Square
type MoveChoiceFunc = func(*Piece, board.Square, board.Square) *Choice
type MoveActionFunc = func(*Piece, board.Square, board.Square, []Option) error

type chainMotions []Motion

func (g chainMotions) Generate(piece *Piece) []*MoveGroup {
	resultMap := make(map[brd.Square]*MoveGroup, 0)
	for _, motion := range g {
		name := motion.Name
		source := piece.Square()
		destinations := motion.MoveGenerator(piece)

		for _, destination := range destinations {
			var optionTree *OptionNode
			if motion.ChoiceFunc != nil {
				choice := motion.ChoiceFunc(piece, source, destination)
				optionTree = choice.GenerateOptions()
			}
			resultMap[destination] = &MoveGroup{
				SquareVec: SquareVec{
					From: piece.Square(),
					To:   destination,
				},
				Name:       name,
				Piece:      piece,
				action:     motion.Action,
				optionTree: optionTree,
			}
		}
	}

	return maps.Values(resultMap)
}

type MoveGroup struct {
	SquareVec
	Name       string
	Piece      *Piece
	action     MoveActionFunc
	optionTree *OptionNode
}

func (mg *MoveGroup) OptionTree() *OptionNode {
	return mg.optionTree
}

func (mg *MoveGroup) Moves() (result []*Move) {
	for _, route := range mg.optionTree.AllRoutes() {
		result = append(result, mg.Move(route))
	}
	return
}

func (mg *MoveGroup) Single() *Move {
	moves := mg.Moves()
	if len(moves) != 1 {
		err := fmt.Errorf("expected move group length of 1, got: %v", len(moves))
		panic(err)
	}
	return moves[0]
}

func (mg *MoveGroup) FilterMoves(predicate func(*Move) bool) *MoveGroup {
	newOptionTree := mg.optionTree.FilterRoutes(func(options []Option) bool {
		move := mg.Move(options)
		return predicate(move)
	})
	return &MoveGroup{
		SquareVec:  mg.SquareVec,
		Name:       mg.Name,
		Piece:      mg.Piece,
		action:     mg.action,
		optionTree: newOptionTree,
	}
}

func (mg *MoveGroup) Move(options []Option) *Move {
	return &Move{
		SquareVec: mg.SquareVec,
		Name:      mg.Name,
		Piece:     mg.Piece,
		Options:   options,
		action:    mg.action,
	}
}

func (mg *MoveGroup) String() string {
	return fmt.Sprintf("%v->%v", mg.From, mg.To)
}

type Move struct {
	SquareVec
	Name    string
	Piece   *Piece
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

type SquareVec struct {
	From brd.Square
	To   brd.Square
}

func (v SquareVec) String() string {
	return fmt.Sprintf("%v->%v", v.From, v.To)
}
