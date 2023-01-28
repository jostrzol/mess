package config

import (
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/userfunc"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/mitchellh/mapstructure"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

type Config struct {
	Board        BoardConfig        `hcl:"board,block"`
	Pieces       PiecesConfig       `hcl:"pieces,block"`
	InitialState InitialStateConfig `hcl:"initial_state,block"`
	Functions    FunctionsConfig
}

type BoardConfig struct {
	Height uint `hcl:"height"`
	Width  uint `hcl:"width"`
}

type PiecesConfig struct {
	Pieces []PieceConfig `hcl:"piece,block"`
}

type PieceConfig struct {
	Name string `hcl:"piece_name,label"`
}

type InitialStateConfig struct {
	PiecePlacements []PiecePlacementConfig `hcl:"piece_placements,block"`
}

type PiecePlacementConfig struct {
	PlayerName string         `hcl:"player_name,label"`
	Placements hcl.Attributes `hcl:",remain"`
}

type FunctionsConfig struct {
	DecideWinner function.Function `mapstructure:"decide_winner"`
}

func ParseFile(filename string) (*Config, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	file, diags := hclsyntax.ParseConfig(src, filename, hcl.InitialPos)
	if diags.HasErrors() {
		return nil, diags
	}

	funcs, body, diags := userfunc.DecodeUserFunctions(file.Body, "function", nil)
	if diags.HasErrors() {
		return nil, diags
	}
	ctx := &hcl.EvalContext{
		Functions: funcs,
	}

	config := &Config{}
	diags = gohcl.DecodeBody(body, ctx, config)
	if diags.HasErrors() {
		return nil, diags
	}

	game := cty.ObjectVal(map[string]cty.Value{
		"players": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("white"),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("black"),
			}),
		}),
	})

	winner, err := funcs["decide_winner"].Call([]cty.Value{game})
	if err != nil {
		return nil, err
	}

	print(winner.GetAttr("name").AsString())

	err = mapstructure.Decode(funcs, &config.Functions)
	if err != nil {
		return nil, err
	}

	return config, nil
}
