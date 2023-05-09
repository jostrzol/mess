package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/jostrzol/mess/config/ctymess"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

var InitialEvalContext = &hcl.EvalContext{
	Functions: map[string]function.Function{
		"upper":               stdlib.UpperFunc,
		"lower":               stdlib.LowerFunc,
		"min":                 stdlib.MinFunc,
		"max":                 stdlib.MaxFunc,
		"strlen":              stdlib.StrlenFunc,
		"substr":              stdlib.SubstrFunc,
		"lookup":              stdlib.LookupFunc,
		"keys":                stdlib.KeysFunc,
		"values":              stdlib.ValuesFunc,
		"element":             stdlib.ElementFunc,
		"length":              stdlib.LengthFunc,
		"range":               stdlib.RangeFunc,
		"slice":               stdlib.SliceFunc,
		"abs":                 stdlib.AbsoluteFunc,
		"sum":                 ctymess.SumFunc,
		"concat":              ctymess.ConcatFunc,
		"all":                 ctymess.AllFunc,
		"square_to_coords":    ctymess.SquareToCoordsFunc,
		"coords_to_square":    ctymess.CoordsToSquareFunc,
		"get_square_relative": ctymess.StateMissingFunc,
		"piece_at":            ctymess.StateMissingFunc,
		"owner_of":            ctymess.StateMissingFunc,
		"is_attacked_by":      ctymess.StateMissingFunc,
		"valid_moves_for":     ctymess.StateMissingFunc,
		"move":                ctymess.StateMissingFunc,
		"capture":             ctymess.StateMissingFunc,
	},
	Variables: map[string]cty.Value{
		"game":  cty.DynamicVal,
		"board": cty.DynamicVal,
	},
}

func populateContextWithState(ctx *hcl.EvalContext, state *mess.State) {
	ctx.Functions["get_square_relative"] = ctymess.GetSquareRelativeFunc(state)
	ctx.Functions["piece_at"] = ctymess.PieceAtFunc(state)
	ctx.Functions["owner_of"] = ctymess.OwnerOfFunc(state)
	ctx.Functions["is_attacked_by"] = ctymess.IsAttackedByFunc(state)
	ctx.Functions["valid_moves_for"] = ctymess.ValidMovesForFunc(state)
	ctx.Functions["move"] = ctymess.MoveFunc(state)
	ctx.Functions["capture"] = ctymess.CaptureFunc(state)

	ctx.Variables["state"] = ctymess.StateToCty(state)
	ctx.Variables["board"] = ctymess.BoardToCty(state.Board())
}
