package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/jostrzol/mess/pkg/rules/ctymess"
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
		"contains":            stdlib.ContainsFunc,
		"format":              stdlib.FormatFunc,
		"sum":                 ctymess.SumFunc,
		"concat":              ctymess.ConcatFunc,
		"all":                 ctymess.AllFunc,
		"any":                 ctymess.AnyFunc,
		"square_to_coords":    ctymess.SquareToCoordsFunc,
		"coords_to_square":    ctymess.CoordsToSquareFunc,
		"println":             ctymess.PrintlnFunc,
		"filternulls":         ctymess.FilterNulls,
		"get_square_relative": ctymess.StateMissingFunc,
		"piece_at":            ctymess.StateMissingFunc,
		"owner_of":            ctymess.StateMissingFunc,
		"is_attacked_by":      ctymess.StateMissingFunc,
		"valid_moves_for":     ctymess.StateMissingFunc,
		"move":                ctymess.StateMissingFunc,
		"capture":             ctymess.StateMissingFunc,
		"place_new_piece":     ctymess.StateMissingFunc,
		"convert_and_release": ctymess.StateMissingFunc,
		"cond_call":           ctymess.StateMissingFunc,
		"call":                ctymess.StateMissingFunc,
	},
	Variables: map[string]cty.Value{
		"game":        cty.DynamicVal,
		"piece_types": cty.DynamicVal,
		"board":       cty.DynamicVal,
	},
}

func initializeContext(ctx *hcl.EvalContext, game *mess.Game) {
	ctx.Functions["get_square_relative"] = ctymess.GetSquareRelativeFunc(game.State)
	ctx.Functions["piece_at"] = ctymess.PieceAtFunc(game.State)
	ctx.Functions["owner_of"] = ctymess.OwnerOfFunc(game.State)
	ctx.Functions["is_attacked_by"] = ctymess.IsAttackedByFunc(game.State)
	ctx.Functions["valid_moves_for"] = ctymess.ValidMovesForFunc(game.State)
	ctx.Functions["move"] = ctymess.MoveFunc(game.State)
	ctx.Functions["capture"] = ctymess.CaptureFunc(game.State)
	ctx.Functions["place_new_piece"] = ctymess.PlaceNewPieceFunc(game.State)
	ctx.Functions["convert_and_release"] = ctymess.ConvertAndReleaseFunc(game.State)
	ctx.Functions["cond_call"] = ctymess.CondCallFunc(ctx)
	ctx.Functions["call"] = ctymess.CallFunc(ctx)

	ctx.Variables["game"] = ctymess.StateToCty(game.State)
	ctx.Variables["piece_types"] = ctymess.PieceTypesToCty(game.PieceTypes())
	ctx.Variables["board"] = ctymess.BoardToCty(game.State.Board())
}
