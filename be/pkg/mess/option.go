package mess

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/board"
)

type Choice interface {
	Message() string
	GenerateOptions() []Option
}
type ChoiceMessage string

func (cb ChoiceMessage) Message() string {
	return string(cb)
}

type Option interface {
	Message() string
	Accept(visitor OptionVisitor)
	String() string
}

type ChoiceGenerator = func([]Option) Choice

// Piece type choice

type PieceTypeChoice struct {
	ChoiceMessage
	PieceTypes []*PieceType
}

func (ptc *PieceTypeChoice) GenerateOptions() []Option {
	result := make([]Option, len(ptc.PieceTypes))
	for i, pieceType := range ptc.PieceTypes {
		result[i] = &PieceTypeOption{ptc.ChoiceMessage, pieceType}
	}
	return result
}

type PieceTypeOption struct {
	ChoiceMessage
	PieceType *PieceType
}

func (pto *PieceTypeOption) Accept(visitor OptionVisitor) {
	visitor.VisitPieceTypeOption(pto)
}

func (pto *PieceTypeOption) String() string {
	return pto.PieceType.name
}

// Square choice

type SquareChoice struct {
	ChoiceMessage
	Squares []board.Square
}

func (sc *SquareChoice) GenerateOptions() []Option {
	result := make([]Option, len(sc.Squares))
	for i, square := range sc.Squares {
		result[i] = &SquareOption{sc.ChoiceMessage, square}
	}
	return result
}

type SquareOption struct {
	ChoiceMessage
	Square board.Square
}

func (so *SquareOption) Accept(visitor OptionVisitor) {
	visitor.VisitSquareOption(so)
}

func (so *SquareOption) String() string {
	return so.Square.String()
}

// Move choice

type MoveChoice struct {
	ChoiceMessage
}

func (mc *MoveChoice) GenerateOptions() []Option {
	return []Option{&MoveOption{ChoiceMessage: mc.ChoiceMessage}}
}

type MoveOption struct {
	ChoiceMessage
	Move *Move
}

func (mo *MoveOption) Accept(visitor OptionVisitor) {
	visitor.VisitMoveOption(mo)
}

func (mo *MoveOption) String() string {
	if mo.Move == nil {
		return "Move option (undecided)"
	}
	return fmt.Sprintf("Move option: %v", mo.Move)
}

// Composite choice

type CompositeChoice struct {
	ChoiceMessage
	Choices []Choice
}

func (cc *CompositeChoice) GenerateOptions() []Option {
	var result []Option
	for _, choice := range cc.Choices {
		result = append(result, choice.GenerateOptions()...)
	}
	return result
}

// Unit choice

type UnitChoice struct {
	ChoiceMessage
}

func (uc *UnitChoice) GenerateOptions() []Option {
	return []Option{&UnitOption{ChoiceMessage: uc.ChoiceMessage}}
}

type UnitOption struct {
	ChoiceMessage
}

func (uo *UnitOption) Accept(visitor OptionVisitor) {
	visitor.VisitUnitOption(uo)
}

func (uo *UnitOption) String() string {
	return "Unit option"
}

// Option visitors

type OptionVisitor interface {
	VisitPieceTypeOption(option *PieceTypeOption)
	VisitSquareOption(option *SquareOption)
	VisitMoveOption(option *MoveOption)
	VisitUnitOption(option *UnitOption)
}
