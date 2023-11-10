package schema

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jostrzol/mess/pkg/mess"
	"github.com/mitchellh/mapstructure"
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
	Option   Option
	Children []*OptionNode
}

type PieceTypeOption PieceType
type SquareOption Square
type MoveOption SquareVec
type UnitOption struct{}

func (o PieceTypeOption) ToDomain(state *mess.State) (mess.Option, error) {
	pieceType, err := state.GetPieceType(o.Name)
	if err != nil {
		return nil, err
	}
	return mess.PieceTypeOption{PieceType: pieceType}, nil
}

func (o SquareOption) ToDomain(_ *mess.State) (mess.Option, error) {
	return mess.SquareOption{Square: Square(o).ToDomain()}, nil
}

func (o MoveOption) ToDomain(_ *mess.State) (mess.Option, error) {
	return mess.MoveOption{SquareVec: SquareVec(o).ToDomain()}, nil
}

func (o UnitOption) ToDomain(_ *mess.State) (mess.Option, error) {
	return mess.UnitOption{}, nil
}

type Option interface {
	ToDomain(state *mess.State) (mess.Option, error)
}

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
		vec := squareVecFromDomain(datum.Option.SquareVec)
		children := optionNodesFromDomain(datum.Children)
		dataMarshalled = append(dataMarshalled, OptionNodeDatum{
			Option:   MoveOption(vec),
			Children: children,
		})
	}
	o.result = &OptionNode{Type: "Move", Message: message, Data: dataMarshalled}
}

func (o *optionTreeMarshaler) VisitUnitNodeData(message string, data mess.UnitOptionNodeData) {
	var dataMarshalled []OptionNodeDatum
	for _, datum := range data.OptionNodeData {
		children := optionNodesFromDomain(datum.Children)
		dataMarshalled = append(dataMarshalled, OptionNodeDatum{
			Option:   UnitOption{},
			Children: children},
		)
	}
	o.result = &OptionNode{Type: "Unit", Message: message, Data: dataMarshalled}
}

type Route []Option

func (r Route) ToDomain(state *mess.State) (result mess.Route, err error) {
	for _, optionDto := range r {
		var option mess.Option
		option, err = optionDto.ToDomain(state)
		result = append(result, option)
	}
	return
}

type RouteBinding struct{}

func (RouteBinding) Name() string {
	return "RouteBinding"
}

func (r RouteBinding) Bind(req *http.Request, obj any) error {
	result, ok := obj.(*Route)
	if !ok {
		return fmt.Errorf("%v excpects a *Route object to bind to, not %T", r.Name(), obj)
	}

	var parsedSlice []any
	bytes, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &parsedSlice)
	if err != nil {
		return err
	}

	for _, parsed := range parsedSlice {
		var tmpOption struct {
			Type string
			Rest map[string]any `mapstructure:",remain"`
		}
		err = mapstructure.Decode(parsed, &tmpOption)
		if err != nil {
			return err
		}
		var option Option
		switch tmpOption.Type {
		case "PieceType":
			option, err = decodeOpton[PieceTypeOption](tmpOption.Rest)
		case "Square":
			option, err = decodeOpton[SquareOption](tmpOption.Rest)
		case "Move":
			option, err = decodeOpton[MoveOption](tmpOption.Rest)
		case "Unit":
			option, err = decodeOpton[UnitOption](tmpOption.Rest)
		}
		if err != nil {
			return err
		}
		*result = append(*result, option)
	}
	return nil
}

func decodeOpton[T any](obj any) (result T, err error) {
	err = mapstructure.Decode(obj, &result)
	return
}
