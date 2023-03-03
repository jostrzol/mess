package config

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/userfunc"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/jostrzol/mess/config/composeuserfunc"
	"github.com/jostrzol/mess/config/messfuncs"
	"github.com/jostrzol/mess/game"
	"github.com/mitchellh/mapstructure"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

func newEvalContext(state *game.State) *hcl.EvalContext {
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
			"sum":                 messfuncs.SumFunc,
			"get_square_relative": messfuncs.GetSquareRelativeFunc(state),
			"is_square_owned_by":  messfuncs.IsSquareOwnedByFunc(state),
			"is_attacked":         messfuncs.IsAttackedFunc(state),
		},
		Variables: make(map[string]cty.Value, 0),
	}
}

type config struct {
	Board        boardConfig        `hcl:"board,block"`
	PieceTypes   pieceTypesConfig   `hcl:"piece_types,block"`
	InitialState initialStateConfig `hcl:"initial_state,block"`
	Functions    *callbackFunctionsConfig
}

type boardConfig struct {
	Height uint `hcl:"height"`
	Width  uint `hcl:"width"`
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
	ActionNames   []string `hcl:"action,optional"`
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

func decodeConfig(filename string, state *game.State) (*config, error) {
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

	ctx := newEvalContext(state)

	userFuncs, body, funcDiags := decodeUserFunctions(file.Body, ctx)
	diags.Extend(funcDiags)

	for name, f := range userFuncs {
		if _, ok := ctx.Functions[name]; ok {
			diags.Append(&hcl.Diagnostic{
				Severity:    hcl.DiagWarning,
				EvalContext: ctx,
				Detail:      fmt.Sprintf("overwrote standard function %q", name),
			})
		}
		ctx.Functions[name] = f
	}

	userVariables, body, varDiags := decodeUserVariables(body, ctx)
	diags.Extend(varDiags)
	ctx.Variables = userVariables

	config := &config{}
	configDiags := gohcl.DecodeBody(body, ctx, config)
	diags.Extend(configDiags)

	err = mapstructure.Decode(ctx.Functions, &config.Functions)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Detail:   fmt.Sprintf("populating callback functions: %v", err),
		})
	}

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
