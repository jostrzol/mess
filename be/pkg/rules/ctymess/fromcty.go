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
	value = tupleToList(value)

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
	var squareStr string

	squareCty, err := getAttr(value, "square")
	if err != nil {
		return nil, err
	}
	square, err := SquareFromCty(squareCty)
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
	var pieceTypeName string

	if err = gocty.FromCtyValue(value, &pieceTypeName); err != nil {
		return nil, fmt.Errorf("parsing piece type name: %w", err)
	}
	for _, pt := range state.PieceTypes() {
		if pt.Name() == pieceTypeName {
			return pt, nil
		}
	}
	return nil, fmt.Errorf("piece type of name %q not found", pieceTypeName)
}

func ChoicesFromCty(state *mess.State, value cty.Value) (choices []*mess.Choice, err error) {
	if value.IsNull() {
		return nil, nil
	}

	value = tupleToList(value)
	if err != nil {
		return nil, err
	}

	if !value.Type().IsListType() {
		return nil, fmt.Errorf("expected list type, found %v", value.Type().FriendlyName())
	}

	result := make([]*mess.Choice, 0, value.LengthInt())
	it := value.ElementIterator()
	for it.Next() {
		index, value := it.Element()
		option, err := ChoiceFromCty(state, value)
		if err != nil {
			return nil, fmt.Errorf("converting element [%d]: %w", index.AsBigFloat(), err)
		}
		result = append(result, option)
	}
	return result, nil
}

func ChoiceFromCty(state *mess.State, value cty.Value) (choice *mess.Choice, err error) {
	if value.IsNull() {
		return nil, nil
	}

	choiceType, err := getAttrAsString(value, "type")
	if err != nil {
		return nil, err
	}
	message, err := getAttrAsString(value, "message")
	if err != nil {
		message = ""
	}
	nextChoicesCty, err := getAttr(value, "next_choices")
	if err != nil {
		nextChoicesCty = cty.NullVal(cty.DynamicPseudoType)
	}
	nextChoices, err := ChoicesFromCty(state, nextChoicesCty)
	if err != nil {
		return nil, err
	}

	var choiceGenerator mess.ChoiceGenerator
	switch choiceType {
	case "piece_type":
		pieceTypeNamesCty, err := getAttr(value, "options")
		if err != nil {
			return nil, err
		}
		pieceTypeNamesCty = tupleToList(pieceTypeNamesCty)

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
		choiceGenerator = &mess.PieceTypeChoiceGenerator{PieceTypes: pieceTypes}
	case "square":
		squaresCty, err := getAttr(value, "squares")
		if err != nil {
			return nil, err
		}

		squares, err := SquaresFromCty(squaresCty)
		if err != nil {
			return nil, err
		}

		choiceGenerator = &mess.SquareChoiceGenerator{Squares: squares}
	case "move":
		choiceGenerator = &mess.MoveChoiceGenerator{State: state}
	case "unit":
		choiceGenerator = &mess.UnitChoiceGenerator{}
	default:
		return nil, fmt.Errorf("invalid choice type: %v", choiceType)
	}
	return &mess.Choice{
		Message:     message,
		NextChoices: nextChoices,
		Generator:   choiceGenerator,
	}, nil
}

func OptionsFromCty(state *mess.State, value cty.Value) (options []mess.Option, err error) {
	if value.IsNull() {
		return nil, nil
	}

	value = tupleToList(value)
	if err != nil {
		return nil, err
	}

	if !value.Type().IsListType() {
		return nil, fmt.Errorf("expected list type, found %v", value.Type().FriendlyName())
	}

	result := make([]mess.Option, 0, value.LengthInt())
	it := value.ElementIterator()
	for it.Next() {
		index, value := it.Element()
		option, err := OptionFromCty(state, value)
		if err != nil {
			return nil, fmt.Errorf("converting element [%d]: %w", index.AsBigFloat(), err)
		}
		result = append(result, option)
	}
	return result, nil
}

func OptionFromCty(state *mess.State, value cty.Value) (mess.Option, error) {
	optionType, err := getAttrAsString(value, "type")
	if err != nil {
		return nil, err
	}

	switch optionType {
	case "piece_type":
		pieceTypeCty, err := getAttr(value, "piece_type", "name")
		if err != nil {
			return nil, err
		}

		pieceType, err := PieceTypeFromCty(state, pieceTypeCty)
		if err != nil {
			return nil, err
		}
		return mess.PieceTypeOption{PieceType: pieceType}, nil
	case "square":
		squareCty, err := getAttr(value, "square")
		if err != nil {
			return nil, err
		}

		square, err := SquareFromCty(squareCty)
		if err != nil {
			return nil, err
		}
		return mess.SquareOption{Square: square}, nil
	case "move":
		moveCty, err := getAttr(value, "move")
		if err != nil {
			return nil, err
		}

		moveGroup, err := MoveGroupFromCty(state, moveCty)
		if err != nil {
			return nil, err
		}
		return mess.MoveOption{MoveGroup: moveGroup}, nil
	case "unit":
		return mess.UnitOption{}, nil
	default:
		return nil, fmt.Errorf("invalid choice type: %v", optionType)
	}
}

func MoveGroupFromCty(state *mess.State, value cty.Value) (*mess.MoveGroup, error) {
	srcCty, err := getAttr(value, "src")
	if err != nil {
		return nil, err
	}
	src, err := SquareFromCty(srcCty)
	if err != nil {
		return nil, err
	}

	dstCty, err := getAttr(value, "dst")
	if err != nil {
		return nil, err
	}
	dst, err := SquareFromCty(dstCty)
	if err != nil {
		return nil, err
	}

	for _, mg := range state.ValidMoves() {
		if mg.From == src && mg.To == dst {
			return mg, nil
		}
	}
	return nil, fmt.Errorf("move not valid")
}

func MoveFromCty(state *mess.State, value cty.Value) (*mess.Move, error) {
	srcCty, err := getAttr(value, "src")
	if err != nil {
		return nil, err
	}
	src, err := SquareFromCty(srcCty)
	if err != nil {
		return nil, err
	}

	dstCty, err := getAttr(value, "dst")
	if err != nil {
		return nil, err
	}
	dst, err := SquareFromCty(dstCty)
	if err != nil {
		return nil, err
	}

	optionsCty, err := getAttr(value, "options")
	if err != nil {
		return nil, err
	}
	options, err := OptionsFromCty(state, optionsCty)
	if err != nil {
		return nil, err
	}

	for _, mg := range state.ValidMoves() {
		if mg.From == src && mg.To == dst {
			return mg.Move(options), nil
		}
	}
	return nil, fmt.Errorf("move not valid")
}

func tupleToList(value cty.Value) cty.Value {
	var err error
	if value.Type().IsTupleType() {
		value, err = convert.Convert(value, cty.List(cty.DynamicPseudoType))
		if err != nil {
			panic(err)
		}
	}
	return value
}

func getAttrAsString(value cty.Value, names ...string) (result string, err error) {
	value, err = getAttr(value, names...)
	if err != nil {
		return "", err
	}

	err = gocty.FromCtyValue(value, &result)
	return
}

func getAttr(value cty.Value, names ...string) (result cty.Value, err error) {
	i := 0
	defer func() {
		recovered := recover()
		switch v := recovered.(type) {
		case nil:
			return
		case string:
			err = fmt.Errorf(v)
		case error:
			err = v
		default:
			panic(v)
		}
		err = fmt.Errorf("getting attribute %q: %w", names[i], err)
	}()

	result = value

	for ; i < len(names); i++ {
		result = result.GetAttr(names[i])
	}
	return
}
