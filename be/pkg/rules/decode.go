package rules

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/userfunc"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/jostrzol/mess/pkg/rules/composeuserfunc"
	"github.com/mitchellh/mapstructure"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

type rules struct {
	Board           boardRules           `hcl:"board,block"`
	PieceTypes      pieceTypesRules      `hcl:"piece_types,block"`
	InitialState    initialStateRules    `hcl:"initial_state,block"`
	StateValidators *stateValidatorRules `hcl:"state_validators,block"`
	Functions       callbackFunctionsRules
}

type boardRules struct {
	Height uint `hcl:"height"`
	Width  uint `hcl:"width"`
}

type pieceTypesRules struct {
	PieceTypes []pieceTypeRules `hcl:"piece_type,block"`
}

type pieceTypeRules struct {
	Name    string        `hcl:"piece_name,label"`
	Symbols *symbols      `hcl:"symbols,block"`
	Motions []motionRules `hcl:"motion,block"`
}

type symbols struct {
	White string `hcl:"white"`
	Black string `hcl:"black"`
}

type motionRules struct {
	GeneratorName string   `hcl:"generator"`
	ActionNames   []string `hcl:"actions,optional"`
}

type initialStateRules struct {
	WhitePieces map[string]string `hcl:"white_pieces"`
	BlackPieces map[string]string `hcl:"black_pieces"`
}

type variablesRules struct {
	VariablesBlock *variablesBlockRules `hcl:"variables,block"`
	Remain         hcl.Body             `hcl:",remain"`
}

type variablesBlockRules struct {
	Variables hcl.Attributes `hcl:",remain"`
}

type stateValidatorRules struct {
	Body hcl.Body `hcl:",remain"`
}

type callbackFunctionsRules struct {
	PickWinnerFunc  function.Function            `mapstructure:"pick_winner"`
	TurnFunc        function.Function            `mapstructure:"turn"`
	CustomFuncs     map[string]function.Function `mapstructure:",remain"`
	StateValidators map[string]function.Function
}

func decodeRules(filename string, ctx *hcl.EvalContext) (*rules, error) {
	diags := make(hcl.Diagnostics, 0)

	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	file, parseDiags := hclsyntax.ParseConfig(src, filename, hcl.InitialPos)
	diags = diags.Extend(parseDiags)
	if diags.HasErrors() {
		return nil, diags
	}

	userFuncs, body, tmpDiags := decodeUserFunctions(file.Body, ctx)
	diags = diags.Extend(tmpDiags)
	tmpDiags = mergeWithStd(ctx.Functions, userFuncs, "function")
	diags = diags.Extend(tmpDiags)

	userVariables, body, tmpDiags := decodeUserVariables(body, ctx)
	diags = diags.Extend(tmpDiags)
	tmpDiags = mergeWithStd(ctx.Variables, userVariables, "variable")
	diags = diags.Extend(tmpDiags)

	rules := &rules{}
	tmpDiags = gohcl.DecodeBody(body, ctx, rules)
	diags = diags.Extend(tmpDiags)

	err = mapstructure.Decode(ctx.Functions, &rules.Functions)
	if err != nil {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Detail:   fmt.Sprintf("populating callback functions: %v", err),
		})
	}

	if rules.StateValidators != nil {
		stateValidators, _, tmpDiags := decodeUserFunctions(rules.StateValidators.Body, ctx)
		diags = diags.Extend(tmpDiags)
		rules.Functions.StateValidators = stateValidators
	} else {
		rules.Functions.StateValidators = make(map[string]function.Function)
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return rules, nil
}

func decodeUserFunctions(
	body hcl.Body, ctx *hcl.EvalContext,
) (map[string]function.Function, hcl.Body, hcl.Diagnostics) {
	diags := make(hcl.Diagnostics, 0)

	contextFunc := func() *hcl.EvalContext {
		return ctx
	}

	userFuncs, remain, tmpDiags := userfunc.DecodeUserFunctions(body, "function", contextFunc)
	diags = diags.Extend(tmpDiags)

	if userFuncs == nil {
		userFuncs = make(map[string]function.Function)
	}

	compositeFuncs, remain, tmpDiags := composeuserfunc.DecodeCompositeUserFunctions(remain, "composite_function", contextFunc)
	diags = diags.Extend(tmpDiags)

	if compositeFuncs == nil {
		compositeFuncs = make(map[string]function.Function)
	}

	for name, f := range compositeFuncs {
		if _, present := userFuncs[name]; present {
			diags = diags.Append(&hcl.Diagnostic{
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
	var variablesRules variablesRules

	tmpDiags := gohcl.DecodeBody(body, ctx, &variablesRules)
	diags = diags.Extend(tmpDiags)

	userVariables := make(map[string]cty.Value)

	if variablesRules.VariablesBlock == nil {
		return userVariables, variablesRules.Remain, diags
	}

	for _, variable := range variablesRules.VariablesBlock.Variables {
		if _, ok := userVariables[variable.Name]; ok {
			diags = diags.Append(&hcl.Diagnostic{
				Severity:    hcl.DiagError,
				Subject:     variable.Expr.Range().Ptr(),
				Expression:  variable.Expr,
				EvalContext: ctx,
				Detail:      fmt.Sprintf("variable named %q already defined", variable.Name),
			})
		} else {
			value, evalDiags := variable.Expr.Value(ctx)
			diags = diags.Extend(evalDiags)
			userVariables[variable.Name] = value
		}
	}

	return userVariables, variablesRules.Remain, diags
}

func mergeWithStd[V any](stdMap map[string]V, userMap map[string]V, kind string) hcl.Diagnostics {
	diags := make(hcl.Diagnostics, 0)
	for name, f := range userMap {
		if _, ok := stdMap[name]; ok {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Detail:   fmt.Sprintf("overwrote standard %s %q", kind, name),
			})
		}
		stdMap[name] = f
	}
	return diags
}
