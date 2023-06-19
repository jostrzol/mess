package config

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/userfunc"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/jostrzol/mess/config/composeuserfunc"
	"github.com/mitchellh/mapstructure"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

type config struct {
	Board           boardConfig          `hcl:"board,block"`
	PieceTypes      pieceTypesConfig     `hcl:"piece_types,block"`
	InitialState    initialStateConfig   `hcl:"initial_state,block"`
	StateValidators stateValidatorConfig `hcl:"state_validators,block"`
	Functions       callbackFunctionsConfig
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
	ActionNames   []string `hcl:"actions,optional"`
}

type initialStateConfig struct {
	WhitePieces map[string]string `hcl:"white_pieces"`
	BlackPieces map[string]string `hcl:"black_pieces"`
}

type variablesConfig struct {
	VariablesBlock *variablesBlockConfig `hcl:"variables,block"`
	Remain         hcl.Body              `hcl:",remain"`
}

type variablesBlockConfig struct {
	Variables hcl.Attributes `hcl:",remain"`
}

type stateValidatorConfig struct {
	Body hcl.Body `hcl:",remain"`
}

type callbackFunctionsConfig struct {
	PickWinnerFunc  function.Function            `mapstructure:"pick_winner"`
	CustomFuncs     map[string]function.Function `mapstructure:",remain"`
	StateValidators map[string]function.Function
}

func decodeConfig(filename string, ctx *hcl.EvalContext) (*config, error) {
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

	userFuncs, body, tmpDiags := decodeUserFunctions(file.Body, ctx)
	diags.Extend(tmpDiags)
	mergeWithStd(ctx.Functions, userFuncs, "function")

	userVariables, body, tmpDiags := decodeUserVariables(body, ctx)
	diags.Extend(tmpDiags)
	mergeWithStd(ctx.Variables, userVariables, "variable")

	config := &config{}
	tmpDiags = gohcl.DecodeBody(body, ctx, config)
	diags.Extend(tmpDiags)

	err = mapstructure.Decode(ctx.Functions, &config.Functions)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Detail:   fmt.Sprintf("populating callback functions: %v", err),
		})
	}

	if config.StateValidators.Body != nil {
		stateValidators, _, tmpDiags := decodeUserFunctions(config.StateValidators.Body, ctx)
		diags.Extend(tmpDiags)
		config.Functions.StateValidators = stateValidators
	} else {
		config.Functions.StateValidators = make(map[string]function.Function)
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
	var variablesConfig variablesConfig

	tmpDiags := gohcl.DecodeBody(body, ctx, &variablesConfig)
	diags.Extend(tmpDiags)

	userVariables := make(map[string]cty.Value)

	if variablesConfig.VariablesBlock == nil {
		return userVariables, variablesConfig.Remain, diags
	}

	for _, variable := range variablesConfig.VariablesBlock.Variables {
		if _, ok := userVariables[variable.Name]; ok {
			diags.Append(&hcl.Diagnostic{
				Severity:    hcl.DiagError,
				Subject:     variable.Expr.Range().Ptr(),
				Expression:  variable.Expr,
				EvalContext: ctx,
				Detail:      fmt.Sprintf("variable named %q already defined", variable.Name),
			})
		} else {
			value, evalDiags := variable.Expr.Value(ctx)
			diags.Extend(evalDiags)
			userVariables[variable.Name] = value
		}
	}

	return userVariables, variablesConfig.Remain, diags
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
