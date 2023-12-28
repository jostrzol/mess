package schema

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jostrzol/mess/pkg/mess"
	"github.com/jostrzol/mess/pkg/server/core/game"
	"github.com/jostrzol/mess/pkg/server/core/usrerr"
	"github.com/mitchellh/mapstructure"
)

func OptionNodeFromDomain(node *mess.OptionNode) *OptionNode {
	if node == nil {
		return nil
	}
	var marshaler optionTreeMarshaler
	node.Accept(&marshaler)
	return marshaler.result
}

func optionNodesFromDomain(nodes []*mess.OptionNode) []*OptionNode {
	result := make([]*OptionNode, 0, len(nodes))
	for _, node := range nodes {
		nodeMarshalled := OptionNodeFromDomain(node)
		if nodeMarshalled != nil {
			result = append(result, nodeMarshalled)
		}
	}
	return result
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

func (o PieceTypeOption) ToDomain(state *game.State) (mess.Option, error) {
	pieceType, ok := state.PieceTypes[o.Name]
	if !ok {
		return nil, usrerr.Errorf("piece type %q not found", o.Name)
	}
	return mess.PieceTypeOption{PieceType: pieceType}, nil
}

func (o SquareOption) ToDomain(_ *game.State) (mess.Option, error) {
	return mess.SquareOption{Square: Square(o).ToDomain()}, nil
}

func (o MoveOption) ToDomain(_ *game.State) (mess.Option, error) {
	return mess.MoveOption{SquareVec: SquareVec(o).ToDomain()}, nil
}

func (o UnitOption) ToDomain(_ *game.State) (mess.Option, error) {
	return mess.UnitOption{}, nil
}

type Option interface {
	ToDomain(state *game.State) (mess.Option, error)
}

type optionTreeMarshaler struct {
	result *OptionNode
}

func (o *optionTreeMarshaler) VisitPieceTypeData(message string, data mess.PieceTypeOptionData) {
	dataMarshalled := []OptionNodeDatum{}
	for _, datum := range data.OptionData {
		pieceType := pieceTypeFromDomain(datum.Option.PieceType)
		children := optionNodesFromDomain(datum.Children)
		dataMarshalled = append(dataMarshalled, OptionNodeDatum{
			Option:   PieceTypeOption(pieceType),
			Children: children,
		})
	}
	o.result = &OptionNode{Type: "PieceType", Message: message, Data: dataMarshalled}
}

func (o *optionTreeMarshaler) VisitSquareData(message string, data mess.SquareOptionData) {
	dataMarshalled := []OptionNodeDatum{}
	for _, datum := range data.OptionData {
		square := squareFromDomain(datum.Option.Square)
		children := optionNodesFromDomain(datum.Children)
		dataMarshalled = append(dataMarshalled, OptionNodeDatum{
			Option:   SquareOption(square),
			Children: children,
		})
	}
	o.result = &OptionNode{Type: "Square", Message: message, Data: dataMarshalled}
}

func (o *optionTreeMarshaler) VisitMoveData(message string, data mess.MoveOptionData) {
	dataMarshalled := []OptionNodeDatum{}
	for _, datum := range data.OptionData {
		vec := squareVecFromDomain(datum.Option.SquareVec)
		children := optionNodesFromDomain(datum.Children)
		dataMarshalled = append(dataMarshalled, OptionNodeDatum{
			Option:   MoveOption(vec),
			Children: children,
		})
	}
	o.result = &OptionNode{Type: "Move", Message: message, Data: dataMarshalled}
}

func (o *optionTreeMarshaler) VisitUnitData(message string, data mess.UnitOptionData) {
	dataMarshalled := []OptionNodeDatum{}
	for _, datum := range data.OptionData {
		children := optionNodesFromDomain(datum.Children)
		dataMarshalled = append(dataMarshalled, OptionNodeDatum{
			Option:   UnitOption{},
			Children: children},
		)
	}
	o.result = &OptionNode{Type: "Unit", Message: message, Data: dataMarshalled}
}

type Route []Option

func (r Route) ToDomain(state *game.State) (result mess.Route, err error) {
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
			option, err = decodeOpton[SquareOption](tmpOption.Rest["Square"])
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
