package rules

import (
	"fmt"

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
	Turn            *turnRules           `hcl:"turn,block"`
	Assets          *cty.Value           `hcl:"assets"`
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
	Name         string         `hcl:"piece_name,label"`
	Presentation *presentations `hcl:"presentation,block"`
	Motions      []motionRules  `hcl:"motion,block"`
}

type presentations struct {
	White *presentation `hcl:"white,block"`
	Black *presentation `hcl:"black,block"`
}

type presentation struct {
	Symbol *string `hcl:"symbol"`
	Icon   *string `hcl:"icon"`
	Rotate *bool   `hcl:"rotate"`
}

type motionRules struct {
	GeneratorName      string `hcl:"generator"`
	ChoiceFunctionName string `hcl:"choice,optional"`
	ActionName         string `hcl:"action,optional"`
}

type initialStateRules struct {
	WhitePieces map[string]string `hcl:"white_pieces"`
	BlackPieces map[string]string `hcl:"black_pieces"`
}

type constantsRules struct {
	ConstantsBlock *constantsBlockRules `hcl:"constants,block"`
	Remain         hcl.Body             `hcl:",remain"`
}

type constantsBlockRules struct {
	Constants hcl.Attributes `hcl:",remain"`
}

type stateValidatorRules struct {
	Body hcl.Body `hcl:",remain"`
}

type turnRules struct {
	ChoiceFunctionName string `hcl:"choice,optional"`
	ActionName         string `hcl:"action,optional"`
}

type callbackFunctionsRules struct {
	ResolutionFunc  function.Function            `mapstructure:"resolve"`
	CustomFuncs     map[string]function.Function `mapstructure:",remain"`
	StateValidators map[string]function.Function
}

func decodeRules(src []byte, filename string, ctx *hcl.EvalContext) (*rules, error) {
	diags := make(hcl.Diagnostics, 0)

	file, parseDiags := hclsyntax.ParseConfig(src, filename, hcl.InitialPos)
	diags = diags.Extend(parseDiags)
	if diags.HasErrors() {
		return nil, diags
	}

	userFuncs, body, tmpDiags := decodeUserFunctions(file.Body, ctx)
	diags = diags.Extend(tmpDiags)
	tmpDiags = mergeWithStd(ctx.Functions, userFuncs, "function")
	diags = diags.Extend(tmpDiags)

	userConstants, body, tmpDiags := decodeUserConstants(body, ctx)
	diags = diags.Extend(tmpDiags)
	tmpDiags = mergeWithStd(ctx.Variables, userConstants, "variable")
	diags = diags.Extend(tmpDiags)

	rules := &rules{}
	tmpDiags = gohcl.DecodeBody(body, ctx, rules)
	diags = diags.Extend(tmpDiags)

	err := mapstructure.Decode(ctx.Functions, &rules.Functions)
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

func decodeUserConstants(
	body hcl.Body, ctx *hcl.EvalContext,
) (map[string]cty.Value, hcl.Body, hcl.Diagnostics) {
	diags := make(hcl.Diagnostics, 0)
	var constantsRules constantsRules

	tmpDiags := gohcl.DecodeBody(body, ctx, &constantsRules)
	diags = diags.Extend(tmpDiags)

	userConstants := make(map[string]cty.Value)

	if constantsRules.ConstantsBlock == nil {
		return userConstants, constantsRules.Remain, diags
	}

	for _, constant := range constantsRules.ConstantsBlock.Constants {
		if _, ok := userConstants[constant.Name]; ok {
			diags = diags.Append(&hcl.Diagnostic{
				Severity:    hcl.DiagError,
				Subject:     constant.Expr.Range().Ptr(),
				Expression:  constant.Expr,
				EvalContext: ctx,
				Detail:      fmt.Sprintf("constant named %q already defined", constant.Name),
			})
		} else {
			value, evalDiags := constant.Expr.Value(ctx)
			diags = diags.Extend(evalDiags)
			userConstants[constant.Name] = value
		}
	}

	return userConstants, constantsRules.Remain, diags
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
