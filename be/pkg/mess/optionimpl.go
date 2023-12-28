package mess

import (
	"github.com/jostrzol/mess/pkg/board"
)

// Visitor

type OptionDataVisitor interface {
	VisitPieceTypeData(message string, data PieceTypeOptionData)
	VisitSquareData(message string, data SquareOptionData)
	VisitMoveData(message string, data MoveOptionData)
	VisitUnitData(message string, data UnitOptionData)
}

// Piece type

type PieceTypeChoice struct {
	PieceTypes []*PieceType
}

func (c *PieceTypeChoice) GenerateOptionData() IOptionData {
	result := make(OptionData[PieceTypeOption], len(c.PieceTypes))
	for i, pieceType := range c.PieceTypes {
		result[i] = &OptionDatum[PieceTypeOption]{
			Option:   PieceTypeOption{pieceType},
			Children: nil,
		}
	}
	return PieceTypeOptionData{result}
}

type PieceTypeOptionData struct {
	OptionData[PieceTypeOption]
}

func (n PieceTypeOptionData) accept(message string, visitor OptionDataVisitor) {
	visitor.VisitPieceTypeData(message, n)
}

func (n PieceTypeOptionData) filter(parentRoute Route, predicate func(Route) bool) IOptionData {
	return PieceTypeOptionData{n.OptionData.filter(parentRoute, predicate)}
}

type PieceTypeOption struct {
	PieceType *PieceType
}

func (o PieceTypeOption) String() string {
	return o.PieceType.Name()
}

// Square choice

type SquareChoice struct {
	Squares []board.Square
}

func (c *SquareChoice) GenerateOptionData() IOptionData {
	result := make(OptionData[SquareOption], len(c.Squares))
	for i, square := range c.Squares {
		result[i] = &OptionDatum[SquareOption]{
			Option:   SquareOption{square},
			Children: nil,
		}
	}
	return SquareOptionData{result}
}

type SquareOptionData struct {
	OptionData[SquareOption]
}

func (n SquareOptionData) accept(message string, visitor OptionDataVisitor) {
	visitor.VisitSquareData(message, n)
}

func (n SquareOptionData) filter(parentRoute Route, predicate func(Route) bool) IOptionData {
	return SquareOptionData{n.OptionData.filter(parentRoute, predicate)}
}

type SquareOption struct {
	Square board.Square
}

func (o SquareOption) String() string {
	return o.Square.String()
}

// Move choice

type MoveChoice struct {
	State *State
}

func (c *MoveChoice) GenerateOptionData() IOptionData {
	validMoves := c.State.ValidMoves()
	result := make(OptionData[MoveOption], 0, len(validMoves))
	for _, moveGroup := range validMoves {
		result = append(result, &OptionDatum[MoveOption]{
			Option:   MoveOption{moveGroup.SquareVec},
			Children: []*OptionNode{moveGroup.optionTree},
		})
	}
	return MoveOptionData{result}
}

type MoveOptionData struct{ OptionData[MoveOption] }

func (n MoveOptionData) accept(message string, visitor OptionDataVisitor) {
	visitor.VisitMoveData(message, n)
}

func (n MoveOptionData) filter(parentRoute Route, predicate func(Route) bool) IOptionData {
	return MoveOptionData{n.OptionData.filter(parentRoute, predicate)}
}

type MoveOption struct {
	SquareVec SquareVec
}

func (o MoveOption) String() string {
	return o.SquareVec.String()
}

// Unit choice

type UnitChoice struct {
}

func (c *UnitChoice) GenerateOptionData() IOptionData {
	return UnitOptionData{
		OptionData[UnitOption]{
			&OptionDatum[UnitOption]{
				Option:   UnitOption{},
				Children: nil,
			},
		},
	}
}

type UnitOptionData struct{ OptionData[UnitOption] }

func (n UnitOptionData) accept(message string, visitor OptionDataVisitor) {
	visitor.VisitUnitData(message, n)
}

func (n UnitOptionData) filter(parentRoute Route, predicate func(Route) bool) IOptionData {
	return UnitOptionData{n.OptionData.filter(parentRoute, predicate)}
}

type UnitOption struct{}

func (o UnitOption) String() string {
	return "<unit>"
}
