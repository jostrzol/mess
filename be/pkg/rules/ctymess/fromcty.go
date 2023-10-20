package ctymess

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"github.com/zclconf/go-cty/cty/gocty"
)

func SquareFromCty(value cty.Value) (board.Square, error) {
	var zero board.Square
	var squareStr string
	if err := gocty.FromCtyValue(value, &squareStr); err != nil {
		return zero, err
	}
	square, err := board.NewSquare(squareStr)
	if err != nil {
		return zero, fmt.Errorf("parsing square %q: %w", squareStr, err)
	}
	return square, nil
}

func SquaresFromCty(value cty.Value) ([]board.Square, error) {
	var squareStrs []string
	var err error
	value, err = tupleToList(value)
	if err != nil {
		return nil, err
	}

	if err := gocty.FromCtyValue(value, &squareStrs); err != nil {
		return nil, err
	}

	errors := make(manyErrors, 0)
	result := make([]board.Square, 0, len(squareStrs))
	for _, squareStr := range squareStrs {
		square, err := board.NewSquare(squareStr)
		if err != nil {
			errors = append(errors, fmt.Errorf("parsing square %q: %w", squareStr, err))
		} else {
			result = append(result, square)
		}
	}
	if len(errors) > 0 {
		return result, errors
	}
	return result, nil
}

func ColorFromCty(colorCty cty.Value) (*color.Color, error) {
	if colorCty.IsNull() {
		return nil, nil
	}
	var colorStr string
	var err error
	if err := gocty.FromCtyValue(colorCty, &colorStr); err != nil {
		return nil, err
	}
	color, err := color.ColorString(colorStr)
	if err != nil {
		return nil, err
	}
	return &color, nil
}

func PieceFromCty(state *mess.State, value cty.Value) (*mess.Piece, error) {
	var err error
	defer func() {
		var ok bool
		value := recover()
		if value == nil {
			return
		}
		err, ok = value.(error)
		if !ok {
			panic(value)
		}
	}()

	squareStr := value.GetAttr("square").AsString()
	if err != nil {
		return nil, fmt.Errorf("getting piece's square: %w", err)
	}
	square, err := board.NewSquare(squareStr)
	if err != nil {
		return nil, fmt.Errorf("parsing square %q: %w", squareStr, err)
	}
	piece, err := state.Board().At(square)
	if err != nil {
		return nil, fmt.Errorf("getting piece at %v: %w", square, err)
	}
	return piece, nil
}

func PieceTypeFromCty(state *mess.State, value cty.Value) (*mess.PieceType, error) {
	var err error
	defer func() {
		var ok bool
		value := recover()
		if value == nil {
			return
		}
		err, ok = value.(error)
		if !ok {
			panic(value)
		}
	}()

	pieceTypeName := value.AsString()
	if err != nil {
		return nil, fmt.Errorf("parsing piece type name: %w", err)
	}
	for _, pt := range state.PieceTypes() {
		if pt.Name() == pieceTypeName {
			return pt, nil
		}
	}
	return nil, fmt.Errorf("piece type of name %q not found", pieceTypeName)
}

func ChoiceFromCty(state *mess.State, value cty.Value) (choice mess.Choice, err error) {
	errPrefix := ""
	defer func() {
		errTmp := recover()
		if errTmp != nil {
			err = fmt.Errorf("%v, %w", errPrefix, errTmp.(error))
		}
	}()

	if value.IsNull() {
		return
	}

	errPrefix = "getting choice type"
	choiceType := value.GetAttr("type").AsString()

	switch choiceType {
	case "piece_type":
		errPrefix = "getting piece_type options"
		pieceTypeNamesCty, err := tupleToList(value.GetAttr("options"))
		if err != nil {
			return nil, err
		}

		var pieceTypeNames []string
		if err := gocty.FromCtyValue(pieceTypeNamesCty, &pieceTypeNames); err != nil {
			return nil, fmt.Errorf("parsing piece_type options: %w", err)
		}

		pieceTypes := make([]*mess.PieceType, len(pieceTypeNames))
		for i, pieceTypeName := range pieceTypeNames {
			pieceType, err := state.GetPieceType(pieceTypeName)
			if err != nil {
				return nil, fmt.Errorf("parsing piece_type choice: %w", err)
			}
			pieceTypes[i] = pieceType
		}
		return mess.PieceTypeChoice{PieceTypes: pieceTypes}, nil
	case "square":
		// TODO:
		panic("TODO")
	default:
		return nil, fmt.Errorf("invalid choice type: %v", choiceType)
	}
}

func tupleToList(value cty.Value) (cty.Value, error) {
	var err error
	if value.Type().IsTupleType() {
		value, err = convert.Convert(value, cty.List(cty.DynamicPseudoType))
		if err != nil {
			return value, fmt.Errorf("transforming to list: %v", err)
		}
	}
	return value, nil
}
