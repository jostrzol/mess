package config

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/userfunc"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/jostrzol/mess/config/composeuserfunc"
	"github.com/jostrzol/mess/config/ctymess"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/mitchellh/mapstructure"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

func newEvalContext(state *mess.State) *hcl.EvalContext {
	return &hcl.EvalContext{
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
			"get_square_relative": ctymess.GetSquareRelativeFunc(state),
			"piece_at":            ctymess.PieceAtFunc(state),
			"owner_of":            ctymess.OwnerOfFunc(state),
			"is_attacked_by":      ctymess.IsAttackedByFunc(state),
			"valid_moves_for":     ctymess.ValidMovesForFunc(state),
			"move":                ctymess.MoveFunc(state),
		},
		Variables: map[string]cty.Value{
			"game": cty.DynamicVal,
		},
	}
}

type minimalConfig struct {
	Board  boardConfig `hcl:"board,block"`
	Remain hcl.Body    `hcl:",remain"`
}

type boardConfig struct {
	Height uint `hcl:"height"`
	Width  uint `hcl:"width"`
}

type config struct {
	State           *mess.State
	PieceTypes      pieceTypesConfig     `hcl:"piece_types,block"`
	InitialState    initialStateConfig   `hcl:"initial_state,block"`
	StateValidators stateValidatorConfig `hcl:"state_validators,block"`
	Functions       callbackFunctionsConfig
}

type pieceTypesConfig struct {
	PieceTypes []pieceTypeConfig `hcl:"piece_type,block"`
}

type pieceTypeConfig struct {
	Name    string         `hcl:"piece_name,label"`
	Motions []motionConfig `hcl:"motion,block"`
}

type motionConfig struct {
	GeneratorName string   `hcl:"generator"`
	ActionNames   []string `hcl:"actions,optional"`
}

type initialStateConfig struct {
	Pieces []piecesConfig `hcl:"pieces,block"`
}

type piecesConfig struct {
	PlayerColor string         `hcl:"player_name,label"`
	Placements  hcl.Attributes `hcl:",remain"`
}

type variablesConfig struct {
	Variables []variableConfig `hcl:"variable,block"`
	Remain    hcl.Body         `hcl:",remain"`
}

type variableConfig struct {
	Name       string         `hcl:"name,label"`
	Expression hcl.Expression `hcl:"value"`
}

type stateValidatorConfig struct {
	Body hcl.Body `hcl:",remain"`
}

func decodeConfig(filename string) (*config, error) {
	diags := make(hcl.Diagnostics, 0)

	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	file, parseDiags := hclsyntax.ParseConfig(src, filename, hcl.InitialPos)
	diags.Extend(parseDiags)
	if diags.HasErrors() {
		return nil, diags
	}

	minConfig := &minimalConfig{}
	tmpDiags := gohcl.DecodeBody(file.Body, nil, minConfig)
	diags.Extend(tmpDiags)

	board, err := mess.NewPieceBoard(int(minConfig.Board.Width), int(minConfig.Board.Height))
	if err != nil {
		return nil, fmt.Errorf("creating new board: %w", err)
	}
	state := mess.NewState(board)
	ctx := newEvalContext(state)

	userFuncs, body, tmpDiags := decodeUserFunctions(minConfig.Remain, ctx)
	diags.Extend(tmpDiags)
	mergeWithStd(ctx.Functions, userFuncs, "function")

	userVariables, body, tmpDiags := decodeUserVariables(body, ctx)
	diags.Extend(tmpDiags)
	mergeWithStd(ctx.Variables, userVariables, "variable")

	config := &config{State: state}
	tmpDiags = gohcl.DecodeBody(body, ctx, config)
	diags.Extend(tmpDiags)

	err = mapstructure.Decode(ctx.Functions, &config.Functions)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Detail:   fmt.Sprintf("populating callback functions: %v", err),
		})
	}

	stateValidators, _, tmpDiags := decodeUserFunctions(config.StateValidators.Body, ctx)
	diags.Extend(tmpDiags)
	config.Functions.StateValidators = stateValidators

	config.Functions.State = state
	config.Functions.EvalContext = ctx

	return config, nil
}

func decodeUserFunctions(
	body hcl.Body, ctx *hcl.EvalContext,
) (map[string]function.Function, hcl.Body, hcl.Diagnostics) {
	diags := make(hcl.Diagnostics, 0)

	contextFunc := func() *hcl.EvalContext {
		return ctx
	}

	userFuncs, remain, tmpDiags := userfunc.DecodeUserFunctions(body, "function", contextFunc)
	diags.Extend(tmpDiags)

	if userFuncs == nil {
		userFuncs = make(map[string]function.Function)
	}

	compositeFuncs, remain, tmpDiags := composeuserfunc.DecodeCompositeUserFunctions(remain, "composite_function", contextFunc)
	diags.Extend(tmpDiags)

	if compositeFuncs == nil {
		compositeFuncs = make(map[string]function.Function)
	}

	for name, f := range compositeFuncs {
		if _, present := userFuncs[name]; present {
			diags.Append(&hcl.Diagnostic{
				Severity:    hcl.DiagError,
				EvalContext: ctx,
				Detail:      fmt.Sprintf("function named %q already defined", name),
			})
		} else {
			userFuncs[name] = f
		}
	}

	return userFuncs, remain, diags
}

func decodeUserVariables(
	body hcl.Body, ctx *hcl.EvalContext,
) (map[string]cty.Value, hcl.Body, hcl.Diagnostics) {
	diags := make(hcl.Diagnostics, 0)
	var variables variablesConfig

	tmpDiags := gohcl.DecodeBody(body, ctx, &variables)
	diags.Extend(tmpDiags)

	if variables.Variables == nil {
		variables.Variables = make([]variableConfig, 0)
	}

	userVariables := make(map[string]cty.Value)
	for _, variable := range variables.Variables {
		if _, ok := userVariables[variable.Name]; ok {
			diags.Append(&hcl.Diagnostic{
				Severity:    hcl.DiagError,
				Subject:     variable.Expression.Range().Ptr(),
				Expression:  variable.Expression,
				EvalContext: ctx,
				Detail:      fmt.Sprintf("variable named %q already defined", variable.Name),
			})
		} else {
			value, evalDiags := variable.Expression.Value(ctx)
			diags.Extend(evalDiags)
			userVariables[variable.Name] = value
		}
	}

	return userVariables, variables.Remain, diags
}

func mergeWithStd[V any](stdMap map[string]V, userMap map[string]V, kind string) hcl.Diagnostics {
	diags := make(hcl.Diagnostics, 0)
	for name, f := range userMap {
		if _, ok := stdMap[name]; ok {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Detail:   fmt.Sprintf("overwrote standard %s %q", kind, name),
			})
		}
		stdMap[name] = f
	}
	return diags
}
