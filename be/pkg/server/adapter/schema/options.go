package schema

import (
	"github.com/jostrzol/mess/pkg/mess"
)

func optionNodesFromDomain(nodes []*mess.OptionNode) []*OptionNode {
	result := make([]*OptionNode, 0, len(nodes))
	for _, node := range nodes {
		nodeMarshalled := optionNodeFromDomain(node)
		if nodeMarshalled != nil {
			result = append(result, nodeMarshalled)
		}
	}
	return result
}

func optionNodeFromDomain(node *mess.OptionNode) *OptionNode {
	if node == nil {
		return nil
	}
	var marshaler optionTreeMarshaler
	node.Accept(&marshaler)
	return marshaler.result
}

type OptionNode struct {
	Type    string
	Message string
	Data    []OptionNodeDatum
}

type OptionNodeDatum struct {
	Option   interface{}
	Children []*OptionNode
}

type PieceTypeOption PieceType
type SquareOption Square
type MoveOption MoveGroup

type optionTreeMarshaler struct {
	result *OptionNode
}

func (o *optionTreeMarshaler) VisitPieceTypeNodeData(message string, data mess.PieceTypeOptionNodeData) {
	var dataMarshalled []OptionNodeDatum
	for _, datum := range data.OptionNodeData {
		pieceType := pieceTypeFromDomain(datum.Option.PieceType)
		children := optionNodesFromDomain(datum.Children)
		dataMarshalled = append(dataMarshalled, OptionNodeDatum{
			Option:   PieceTypeOption(pieceType),
			Children: children,
		})
	}
	o.result = &OptionNode{Type: "PieceType", Message: message, Data: dataMarshalled}
}

func (o *optionTreeMarshaler) VisitSquareNodeData(message string, data mess.SquareOptionNodeData) {
	var dataMarshalled []OptionNodeDatum
	for _, datum := range data.OptionNodeData {
		square := squareFromDomain(datum.Option.Square)
		children := optionNodesFromDomain(datum.Children)
		dataMarshalled = append(dataMarshalled, OptionNodeDatum{
			Option:   SquareOption(square),
			Children: children,
		})
	}
	o.result = &OptionNode{Type: "Square", Message: message, Data: dataMarshalled}
}

func (o *optionTreeMarshaler) VisitMoveNodeData(message string, data mess.MoveOptionNodeData) {
	var dataMarshalled []OptionNodeDatum
	for _, datum := range data.OptionNodeData {
		moveGroup := moveGroupFromDomain(datum.Option.MoveGroup)
		children := optionNodesFromDomain(datum.Children)
		dataMarshalled = append(dataMarshalled, OptionNodeDatum{
			Option:   MoveOption(moveGroup),
			Children: children,
		})
	}
	o.result = &OptionNode{Type: "Move", Message: message, Data: dataMarshalled}
}

func (o *optionTreeMarshaler) VisitUnitNodeData(message string, data mess.UnitOptionNodeData) {
	var dataMarshalled []OptionNodeDatum
	for _, datum := range data.OptionNodeData {
		children := optionNodesFromDomain(datum.Children)
		dataMarshalled = append(dataMarshalled, OptionNodeDatum{Children: children})
	}
	o.result = &OptionNode{Type: "Unit", Message: message, Data: dataMarshalled}
}
