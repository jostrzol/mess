package config

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/userfunc"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/jostrzol/mess/config/composeuserfunc"
	"github.com/jostrzol/mess/config/messfuncs"
	"github.com/mitchellh/mapstructure"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

var defaultEvalContext = &hcl.EvalContext{
	Functions: map[string]function.Function{
		"upper":   stdlib.UpperFunc,
		"lower":   stdlib.LowerFunc,
		"min":     stdlib.MinFunc,
		"max":     stdlib.MaxFunc,
		"strlen":  stdlib.StrlenFunc,
		"substr":  stdlib.SubstrFunc,
		"lookup":  stdlib.LookupFunc,
		"keys":    stdlib.KeysFunc,
		"values":  stdlib.ValuesFunc,
		"element": stdlib.ElementFunc,
		"length":  stdlib.LengthFunc,
		"sum":     messfuncs.SumFunc,
	},
	Variables: make(map[string]cty.Value, 0),
}

type config struct {
	Board        boardConfig        `hcl:"board,block"`
	PieceTypes   pieceTypesConfig   `hcl:"piece_types,block"`
	InitialState initialStateConfig `hcl:"initial_state,block"`
	Functions    CallbackFunctionsConfig
}

type boardConfig struct {
	Height uint `hcl:"height"`
	Width  uint `hcl:"width"`
}

type pieceTypesConfig struct {
	PieceTypes []pieceTypeConfig `hcl:"piece_type,block"`
}

type pieceTypeConfig struct {
	Name string `hcl:"piece_name,label"`
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
	Name  string         `hcl:"name,label"`
	Value hcl.Expression `hcl:"value"`
}

func DecodeConfig(filename string) (*config, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	file, diags := hclsyntax.ParseConfig(src, filename, hcl.InitialPos)
	if diags.HasErrors() {
		return nil, diags
	}

	ctx := *defaultEvalContext

	variables := &variablesConfig{}
	diags = gohcl.DecodeBody(file.Body, &ctx, variables)
	if diags.HasErrors() {
		return nil, diags
	}

	for _, v := range variables.Variables {
		if _, ok := defaultEvalContext.Variables[v.Name]; ok {
			log.Printf("user overwrote standard variable %q!", v.Name)
		}

		ctx.Variables[v.Name], diags = v.Value.Value(&ctx)
		if diags.HasErrors() {
			return nil, diags
		}
	}

	contextFunc := func() *hcl.EvalContext {
		return &ctx
	}

	funcs, body, diags := userfunc.DecodeUserFunctions(variables.Remain, "function", contextFunc)
	if diags.HasErrors() {
		return nil, diags
	}

	for name, f := range funcs {
		if _, ok := defaultEvalContext.Functions[name]; ok {
			log.Printf("user overwrote standard function %q!", name)
		}
		ctx.Functions[name] = f
	}

	compositeFuncs, body, diags := composeuserfunc.DecodeCompositeUserFunctions(body, "composite_function", contextFunc)
	if diags.HasErrors() {
		return nil, diags
	}

	for name, f := range compositeFuncs {
		if _, ok := funcs[name]; ok {
			return nil, fmt.Errorf("user function name clash: %q", name)
		} else if _, ok := defaultEvalContext.Functions[name]; ok {
			log.Printf("user overwrote standard function %q!", name)
		}
		ctx.Functions[name] = f
	}

	config := &config{}
	diags = gohcl.DecodeBody(body, &ctx, config)
	if diags.HasErrors() {
		return nil, diags
	}

	err = mapstructure.Decode(ctx.Functions, &config.Functions)
	if err != nil {
		return nil, err
	}

	return config, nil
}
