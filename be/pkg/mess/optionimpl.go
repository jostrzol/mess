package mess

import (
	"github.com/jostrzol/mess/pkg/board"
)

// Visitor

type OptionNodeDataVisitor interface {
	VisitPieceTypeNodeData(message string, data PieceTypeOptionNodeData)
	VisitSquareNodeData(message string, data SquareOptionNodeData)
	VisitMoveNodeData(message string, data MoveOptionNodeData)
	VisitUnitNodeData(message string, data UnitOptionNodeData)
}

// Piece type

type PieceTypeChoice struct {
	PieceTypes []*PieceType
}

func (c *PieceTypeChoice) GenerateOptions() IOptionNodeData {
	result := make(OptionNodeData[PieceTypeOption], len(c.PieceTypes))
	for i, pieceType := range c.PieceTypes {
		result[i] = &OptionNodeDatum[PieceTypeOption]{
			Option:   PieceTypeOption{pieceType},
			Children: nil,
		}
	}
	return PieceTypeOptionNodeData{result}
}

type PieceTypeOptionNodeData struct {
	OptionNodeData[PieceTypeOption]
}

func (n PieceTypeOptionNodeData) accept(message string, visitor OptionNodeDataVisitor) {
	visitor.VisitPieceTypeNodeData(message, n)
}

func (n PieceTypeOptionNodeData) filter(parentRoute Route, predicate func(Route) bool) IOptionNodeData {
	return PieceTypeOptionNodeData{n.OptionNodeData.filter(parentRoute, predicate)}
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

func (c *SquareChoice) GenerateOptions() IOptionNodeData {
	result := make(OptionNodeData[SquareOption], len(c.Squares))
	for i, square := range c.Squares {
		result[i] = &OptionNodeDatum[SquareOption]{
			Option:   SquareOption{square},
			Children: nil,
		}
	}
	return SquareOptionNodeData{result}
}

type SquareOptionNodeData struct {
	OptionNodeData[SquareOption]
}

func (n SquareOptionNodeData) accept(message string, visitor OptionNodeDataVisitor) {
	visitor.VisitSquareNodeData(message, n)
}

func (n SquareOptionNodeData) filter(parentRoute Route, predicate func(Route) bool) IOptionNodeData {
	return SquareOptionNodeData{n.OptionNodeData.filter(parentRoute, predicate)}
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

func (c *MoveChoice) GenerateOptions() IOptionNodeData {
	validMoves := c.State.ValidMoves()
	result := make(OptionNodeData[MoveOption], 0, len(validMoves))
	for _, moveGroup := range validMoves {
		result = append(result, &OptionNodeDatum[MoveOption]{
			Option:   MoveOption{moveGroup.SquareVec},
			Children: []*OptionNode{moveGroup.optionTree},
		})
	}
	return MoveOptionNodeData{result}
}

type MoveOptionNodeData struct{ OptionNodeData[MoveOption] }

func (n MoveOptionNodeData) accept(message string, visitor OptionNodeDataVisitor) {
	visitor.VisitMoveNodeData(message, n)
}

func (n MoveOptionNodeData) filter(parentRoute Route, predicate func(Route) bool) IOptionNodeData {
	return MoveOptionNodeData{n.OptionNodeData.filter(parentRoute, predicate)}
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

func (c *UnitChoice) GenerateOptions() IOptionNodeData {
	return UnitOptionNodeData{
		OptionNodeData[UnitOption]{
			&OptionNodeDatum[UnitOption]{
				Option:   UnitOption{},
				Children: nil,
			},
		},
	}
}

type UnitOptionNodeData struct{ OptionNodeData[UnitOption] }

func (n UnitOptionNodeData) accept(message string, visitor OptionNodeDataVisitor) {
	visitor.VisitUnitNodeData(message, n)
}

func (n UnitOptionNodeData) filter(parentRoute Route, predicate func(Route) bool) IOptionNodeData {
	return UnitOptionNodeData{n.OptionNodeData.filter(parentRoute, predicate)}
}

type UnitOption struct{}

func (o UnitOption) String() string {
	return "<unit>"
}
