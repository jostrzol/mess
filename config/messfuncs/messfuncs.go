package messfuncs

import (
	"fmt"
	"strings"

	"github.com/jostrzol/mess/config/ctyconv"
	"github.com/jostrzol/mess/game"
	"github.com/jostrzol/mess/game/board"
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

func GetSquareRelativeFunc(state *game.State) function.Function {
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
		Type: function.StaticReturnType(cty.Number),
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
