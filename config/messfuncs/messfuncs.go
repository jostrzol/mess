package messfuncs

import (
	"fmt"
	"strings"

	"github.com/jostrzol/mess/config/ctyconv"
	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"
)

func joinText(text ...string) string {
	return strings.Join(text, "\n")
}

var SumFunc = function.New(&function.Spec{
	Description: "Sums all the given numbers",
	Params:      []function.Parameter{},
	VarParam: &function.Parameter{
		Name:             "numbers",
		Type:             cty.Number,
		AllowDynamicType: true,
	},
	Type: function.StaticReturnType(cty.Number),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		sum := cty.Zero
		for _, num := range args {
			sum = sum.Add(num)
		}

		return sum, nil
	},
})

var ctyOffset = cty.Tuple([]cty.Type{cty.Number, cty.Number})

func GetSquareRelativeFunc(state *mess.State) function.Function {
	return function.New(&function.Spec{
		Description: joinText(
			"Gets the square offset by a given relative position,",
			"or nil if the board doesn't contain the square",
		),
		Params: []function.Parameter{
			{
				Name:             "square",
				Type:             cty.String,
				AllowDynamicType: true,
			},
			{
				Name:             "offset",
				Type:             ctyOffset,
				AllowDynamicType: true,
			},
		},
		Type: function.StaticReturnType(cty.String),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			var square *board.Square
			var offset board.Offset
			var err error
			if square, err = ctyconv.SquareFromCty(args[0]); err != nil {
				return cty.DynamicVal, fmt.Errorf("argument 'square': %w", err)
			}
			if err = gocty.FromCtyValue(args[1], &offset); err != nil {
				return cty.DynamicVal, fmt.Errorf("argument 'offset': %w", err)
			}

			result := square.Offset(offset)
			if !state.Board.Contains(result) {
				return cty.NullVal(cty.String), nil
			}
			return ctyconv.SquareToCty(result), nil
		},
	})
}

func IsSquareOwnedByFunc(state *mess.State) function.Function {
	return function.New(&function.Spec{
		Description: "Checks if a square has a piece of a given color",
		Params: []function.Parameter{
			{
				Name:             "square",
				Type:             cty.String,
				AllowDynamicType: true,
				AllowNull:        true,
			},
			{
				Name:             "color",
				Type:             cty.String,
				AllowDynamicType: true,
				AllowNull:        true,
			},
		},
		Type: function.StaticReturnType(cty.Bool),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			if args[0].IsNull() || args[1].IsNull() {
				return cty.BoolVal(false), nil
			}
			var square *board.Square
			var err error
			if square, err = ctyconv.SquareFromCty(args[0]); err != nil {
				return cty.DynamicVal, fmt.Errorf("argument 'square': %w", err)
			}
			color, err := ctyconv.ColorFromCty(args[1])
			if err != nil {
				return cty.DynamicVal, fmt.Errorf("argument 'color': %w", err)
			}

			piece, err := state.PieceAt(square)
			if err != nil {
				return cty.DynamicVal, fmt.Errorf("getting piece at %v: %w", square, err)
			}
			result := piece != nil && piece.Color() == *color
			return cty.BoolVal(result), nil
		},
	})
}

func IsAttackedFunc(state *mess.State) function.Function {
	return function.New(&function.Spec{
		Description: "Checks if given square can be reached in the next turn by the opponent",
		Params: []function.Parameter{
			{
				Name:             "square",
				Type:             cty.String,
				AllowDynamicType: true,
			},
		},
		Type: function.StaticReturnType(cty.Bool),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			// TODO: implement
			return cty.BoolVal(false), nil
		},
	})
}
