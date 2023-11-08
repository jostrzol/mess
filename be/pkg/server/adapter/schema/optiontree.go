package schema

import (
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/jostrzol/mess/pkg/utils"
)

func optionTreeFromDomain(optionTree mess.OptionTree) interface{} {
	if optionTree == nil {
		return nil
	}
	var marshaler optionTreeMarshaler
	optionTree.Accept(&marshaler)
	return marshaler.result
}

type PieceTypeNode struct {
	Type    string
	Options []PieceTypeNodeOption
}

type PieceTypeNodeOption struct {
	PieceType PieceType
	Next      interface{} `json:",omitempty"`
}

type SquareNode struct {
	Type    string
	Options []SquareNodeOption
}

type SquareNodeOption struct {
	Square Square
	Next   interface{} `json:",omitempty"`
}

type MoveNode struct {
	Type string
	Next interface{} `json:",omitempty"`
}

type UnitNode struct {
	Type string
	Next interface{} `json:",omitempty"`
}

type MessageNode struct {
	Type     string
	Children map[string]interface{}
}

type StopActionNode struct {
	Type string
}

type optionTreeMarshaler struct {
	result interface{}
}

func (o *optionTreeMarshaler) VisitPieceTypeNode(options map[*mess.PieceTypeOption]mess.OptionTree) {
	var nodeOptions []PieceTypeNodeOption
	for option, node := range options {
		pieceType := pieceTypeFromDomain(option.PieceType)
		next := optionTreeFromDomain(node)
		nodeOptions = append(nodeOptions, PieceTypeNodeOption{PieceType: pieceType, Next: next})
	}
	o.result = PieceTypeNode{Type: "PieceType", Options: nodeOptions}
}

func (o *optionTreeMarshaler) VisitSquareNode(options map[*mess.SquareOption]mess.OptionTree) {
	var nodeOptions []SquareNodeOption
	for option, node := range options {
		square := squareFromDomain(option.Square)
		next := optionTreeFromDomain(node)
		nodeOptions = append(nodeOptions, SquareNodeOption{Square: square, Next: next})
	}
	o.result = SquareNode{Type: "Square", Options: nodeOptions}
}

func (o *optionTreeMarshaler) VisitMoveNode(options map[*mess.MoveOption]mess.OptionTree) {
	_, node := utils.SingleEntry(options)
	o.result = MoveNode{Type: "Move", Next: optionTreeFromDomain(node)}
}

func (o *optionTreeMarshaler) VisitUnitNode(options map[*mess.UnitOption]mess.OptionTree) {
	_, node := utils.SingleEntry(options)
	o.result = MoveNode{Type: "Unit", Next: optionTreeFromDomain(node)}
}

func (o *optionTreeMarshaler) VisitMessageNode(options map[string]mess.OptionTree) {
	if len(options) == 0 {
		o.result = StopActionNode{Type: "StopAction"}
	} else {
		children := make(map[string]interface{})
		for msg, node := range options {
			node.Accept(o)
			children[msg] = o.result
		}
		o.result = children
	}
}