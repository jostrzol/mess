package ctymess

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/color"
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

var ConcatFunc = function.New(&function.Spec{
	Description: "Concatenates the given lists",
	Params:      []function.Parameter{},
	VarParam: &function.Parameter{
		Name:             "lists",
		Type:             cty.List(cty.DynamicPseudoType),
		AllowDynamicType: true,
	},
	Type: function.StaticReturnType(cty.List(cty.DynamicPseudoType)),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		result := make([]cty.Value, 0)
		for _, list := range args {
			result = append(result, list.AsValueSlice()...)
		}

		return listOrEmpty(cty.DynamicPseudoType, result), nil
	},
})

var AllFunc = function.New(&function.Spec{
	Description: "Returns true only if all arguments are true",
	Params:      []function.Parameter{},
	VarParam: &function.Parameter{
		Name:             "args",
		Type:             cty.Bool,
		AllowDynamicType: true,
	},
	Type: function.StaticReturnType(cty.Bool),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		result := cty.True
		for _, arg := range args {
			result = result.And(arg)
		}

		return result, nil
	},
})

var SquareToCoordsFunc = function.New(&function.Spec{
	Description: "Converts a square in string format to a tuple of numbers {file, rank}",
	Params: []function.Parameter{
		{
			Name:             "square",
			Type:             cty.String,
			AllowDynamicType: true,
		},
	},
	Type: function.StaticReturnType(Coords),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		var square board.Square
		var err error
		if square, err = SquareFromCty(args[0]); err != nil {
			return cty.DynamicVal, fmt.Errorf("argument 'square': %w", err)
		}

		x, y := square.ToCoords()

		return cty.TupleVal(
			[]cty.Value{
				cty.NumberIntVal(int64(x)),
				cty.NumberIntVal(int64(y)),
			}), nil
	},
})

var CoordsToSquareFunc = function.New(&function.Spec{
	Description: "Converts coords in tuple of numbers {file, rank} to a string square",
	Params: []function.Parameter{
		{
			Name:             "coords",
			Type:             Coords,
			AllowDynamicType: true,
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		var coords struct {
			X int
			Y int
		}

		if err := gocty.FromCtyValue(args[0], &coords); err != nil {
			return cty.DynamicVal, fmt.Errorf("argument 'coords': %w", err)
		}

		square := board.SquareFromCoords(coords.X, coords.Y)

		return cty.StringVal(square.String()), nil
	},
})

var RangeFunc = function.New(&function.Spec{
	Description: "Returns a list of numbers in the given range (include start, exclude end)",
	Params: []function.Parameter{
		{
			Name:             "start",
			Type:             cty.Number,
			AllowDynamicType: true,
		},
		{
			Name:             "end",
			Type:             cty.Number,
			AllowDynamicType: true,
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		var coords struct {
			x int
			y int
		}

		if err := gocty.FromCtyValue(args[0], &coords); err != nil {
			return cty.DynamicVal, fmt.Errorf("argument 'coords': %w", err)
		}

		square := board.SquareFromCoords(coords.x, coords.y)

		return cty.StringVal(square.String()), nil
	},
})

var PrintlnFunc = function.New(&function.Spec{
	Description: "Prints its argument",
	Params: []function.Parameter{
		{
			Name:             "argument",
			Type:             cty.DynamicPseudoType,
			AllowDynamicType: true,
		},
	},
	Type: function.StaticReturnType(cty.DynamicPseudoType),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		fmt.Printf("%#v\n", args[0])
		return cty.DynamicVal, nil
	},
})

var StateMissingFunc = function.New(&function.Spec{
	Description: "Placeholder - the game state is not built yet",
	Params:      []function.Parameter{},
	Type:        function.StaticReturnType(cty.NilType),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		panic("state not built yet")
	},
})

func GetSquareRelativeFunc(state *mess.State) function.Function {
	return function.New(&function.Spec{
		Description: joinText(
			"Gets the square offset by a given relative position,",
			"or null if the board doesn't contain the square",
		),
		Params: []function.Parameter{
			{
				Name:             "square",
				Type:             cty.String,
				AllowDynamicType: true,
			},
			{
				Name:             "offset",
				Type:             Offset,
				AllowDynamicType: true,
			},
		},
		Type: function.StaticReturnType(cty.String),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			var square board.Square
			var offset board.Offset
			var err error
			if square, err = SquareFromCty(args[0]); err != nil {
				return cty.DynamicVal, fmt.Errorf("argument 'square': %w", err)
			}
			if err = gocty.FromCtyValue(args[1], &offset); err != nil {
				return cty.DynamicVal, fmt.Errorf("argument 'offset': %w", err)
			}

			result := square.Offset(offset)
			if !state.Board().Contains(result) {
				return cty.NullVal(cty.String), nil
			}
			return SquareToCty(result), nil
		},
	})
}

func PieceAtFunc(state *mess.State) function.Function {
	return function.New(&function.Spec{
		Description: joinText(
			"Get piece at the given square or null if either the board doesn't",
			"contain the square or there is no piece there",
		),
		Params: []function.Parameter{
			{
				Name:             "square",
				Type:             cty.String,
				AllowDynamicType: true,
			},
		},
		Type: function.StaticReturnType(Piece),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			var square board.Square
			var err error
			if square, err = SquareFromCty(args[0]); err != nil {
				return cty.DynamicVal, fmt.Errorf("argument 'square': %w", err)
			}

			piece, err := state.Board().At(square)
			if err != nil {
				log.Printf("getting piece at %v - returning null: %v", square, err)
				return cty.NullVal(Piece), nil
			} else if piece == nil {
				return cty.NullVal(Piece), nil
			}

			return PieceToCty(piece), nil
		},
	})
}

func OwnerOfFunc(state *mess.State) function.Function {
	return function.New(&function.Spec{
		Description: "Get owner of the given piece",
		Params: []function.Parameter{
			{
				Name:             "piece",
				Type:             Piece,
				AllowDynamicType: true,
			},
		},
		Type: function.StaticReturnType(Player),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			squareCty := args[0].GetAttr("square")

			var square board.Square
			var err error
			if square, err = SquareFromCty(squareCty); err != nil {
				return cty.DynamicVal, fmt.Errorf("argument 'square': %w", err)
			}

			piece, err := state.Board().At(square)
			if err != nil {
				return cty.DynamicVal, fmt.Errorf("getting piece at %v: %v", square, err)

			} else if piece == nil {
				return cty.DynamicVal, fmt.Errorf("getting piece at %v: no piece", square)
			}

			return PlayerToCty(piece.Owner()), nil
		},
	})
}

// TODO: block usage in generators (don't include it in EvalContext for generators).
func IsAttackedByFunc(state *mess.State) function.Function {
	return function.New(&function.Spec{
		Description: "Checks if player of the given color can reach the given square",
		Params: []function.Parameter{
			{
				Name:             "color",
				Type:             cty.String,
				AllowDynamicType: true,
			},
			{
				Name:             "square",
				Type:             cty.String,
				AllowDynamicType: true,
			},
		},
		Type: function.StaticReturnType(cty.Bool),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			var square board.Square
			var color *color.Color
			var err error

			if color, err = ColorFromCty(args[0]); err != nil {
				return cty.DynamicVal, fmt.Errorf("argument 'color': %w", err)
			}
			if square, err = SquareFromCty(args[1]); err != nil {
				return cty.DynamicVal, fmt.Errorf("argument 'square': %w", err)
			}

			player := state.Player(*color)
			for _, attacked := range player.AttackedSquares() {
				if attacked == square {
					return cty.True, nil
				}
			}
			return cty.False, nil
		},
	})
}

func ValidMovesForFunc(state *mess.State) function.Function {
	return function.New(&function.Spec{
		Description: "Returns all the squares that the given piece can go to in 1 turn",
		Params: []function.Parameter{
			{
				Name:             "piece",
				Type:             Piece,
				AllowDynamicType: true,
			},
		},
		Type: function.StaticReturnType(cty.List(Move)),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			var piece *mess.Piece
			var err error
			if piece, err = PieceFromCty(state, args[0]); err != nil {
				return cty.DynamicVal, fmt.Errorf("argument 'piece': %w", err)
			} else if piece == nil {
				return cty.DynamicVal, fmt.Errorf("given piece not found")
			}

			result := make([]cty.Value, 0)
			for _, move := range state.ValidMoves() {
				if move.Piece == piece {
					result = append(result, MoveToCty(&move))
				}
			}

			return listOrEmpty(Move, result), nil
		},
	})
}

func MoveFunc(state *mess.State) function.Function {
	return function.New(&function.Spec{
		Description: "Moves the given piece to the given square",
		Params: []function.Parameter{
			{
				Name:             "piece",
				Type:             Piece,
				AllowDynamicType: true,
			},
			{
				Name:             "destination",
				Type:             cty.String,
				AllowDynamicType: true,
			},
		},
		Type: function.StaticReturnType(cty.EmptyTuple),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			var piece *mess.Piece
			var destination board.Square
			var err error
			if piece, err = PieceFromCty(state, args[0]); err != nil {
				return cty.DynamicVal, fmt.Errorf("argument 'piece': %w", err)
			} else if piece == nil {
				return cty.DynamicVal, fmt.Errorf("given piece not found")
			}
			if destination, err = SquareFromCty(args[1]); err != nil {
				return cty.DynamicVal, fmt.Errorf("argument 'destination': %w", err)
			}

			if err = piece.MoveTo(destination); err != nil {
				return cty.DynamicVal, fmt.Errorf("moving %v to %v: %w", piece, destination, err)
			}

			return cty.EmptyTupleVal, nil
		},
	})
}

func CaptureFunc(state *mess.State) function.Function {
	return function.New(&function.Spec{
		Description: "Captures the given piece",
		Params: []function.Parameter{
			{
				Name:             "piece",
				Type:             Piece,
				AllowDynamicType: true,
			},
		},
		Type: function.StaticReturnType(cty.EmptyTuple),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			var piece *mess.Piece
			var err error
			if piece, err = PieceFromCty(state, args[0]); err != nil {
				return cty.DynamicVal, fmt.Errorf("argument 'piece': %w", err)
			} else if piece == nil {
				return cty.DynamicVal, fmt.Errorf("given piece not found")
			}

			opponent := state.OpponentTo(piece.Owner())

			if err = piece.GetCapturedBy(opponent); err != nil {
				return cty.DynamicVal, fmt.Errorf("capturing %v: %w", piece, err)
			}

			return cty.EmptyTupleVal, nil
		},
	})
}

func ChooseFunc(game *mess.Game) function.Function {
	return function.New(&function.Spec{
		Description: "Makes the player choose one of the provided options",
		Params: []function.Parameter{
			{
				Name:             "options",
				Type:             cty.List(cty.String),
				AllowDynamicType: true,
			},
		},
		Type: function.StaticReturnType(cty.Number),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			var options []string

			err := gocty.FromCtyValue(args[0], &options)
			if err != nil {
				return cty.DynamicVal, fmt.Errorf("argument 'options': %w", err)
			}

			if len(options) == 0 {
				return cty.NullVal(cty.String), nil
			}

			if game.IsGeneratingMoves() {
				return cty.StringVal(options[0]), nil
			}

			chosen := game.Choose(options)

			return cty.NumberIntVal(int64(chosen)), nil
		},
	})
}

func PlaceNewPieceFunc(state *mess.State) function.Function {
	return function.New(&function.Spec{
		Description: "Place a new piece at the given destination. Replace if occupied.",
		Params: []function.Parameter{
			{
				Name:             "piece_type_name",
				Type:             cty.String,
				AllowDynamicType: true,
			},
			{
				Name:             "square",
				Type:             cty.String,
				AllowDynamicType: true,
			},
			{
				Name:             "color",
				Type:             cty.String,
				AllowDynamicType: true,
			},
		},
		Type: function.StaticReturnType(cty.DynamicPseudoType),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			var pieceType *mess.PieceType
			var square board.Square
			var color *color.Color
			var err error

			if pieceType, err = PieceTypeFromCty(state, args[0]); err != nil {
				return cty.DynamicVal, fmt.Errorf("argument 'pieceType': %w", err)
			}
			if square, err = SquareFromCty(args[1]); err != nil {
				return cty.DynamicVal, fmt.Errorf("argument 'square': %w", err)
			}
			if color, err = ColorFromCty(args[2]); err != nil {
				return cty.DynamicVal, fmt.Errorf("argument 'color': %w", err)
			}

			piece := mess.NewPiece(pieceType, state.Player(*color))
			if err = state.Board().Replace(piece, square); err != nil {
				return cty.DynamicVal, fmt.Errorf("placing new piece: %w", err)
			}

			return cty.DynamicVal, nil
		},
	})
}

func CondCallFunc(ctx *hcl.EvalContext) function.Function {
	return function.New(&function.Spec{
		Description: "Calls the given function, if the condition is true, else returns null",
		Params: []function.Parameter{
			{
				Name:             "condition",
				Type:             cty.Bool,
				AllowDynamicType: true,
			},
			{
				Name:             "function_name",
				Type:             cty.DynamicPseudoType,
				AllowDynamicType: true,
			},
		},
		VarParam: &function.Parameter{
			Name:             "arguments",
			Type:             cty.DynamicPseudoType,
			AllowDynamicType: true,
		},
		Type: function.StaticReturnType(cty.DynamicPseudoType),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			var err error
			if args[0].True() {
				function := ctx.Functions[args[1].AsString()]
				_, err = function.Call(args[2:])
			}
			return cty.NullVal(cty.DynamicPseudoType), err
		},
	})
}
