package config

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/userfunc"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/jostrzol/mess/game"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	Board        BoardConfig        `hcl:"board,block"`
	PieceTypes   PieceTypesConfig   `hcl:"piece_types,block"`
	InitialState InitialStateConfig `hcl:"initial_state,block"`
	Functions    FunctionsConfig
}

type BoardConfig struct {
	Height uint `hcl:"height"`
	Width  uint `hcl:"width"`
}

type PieceTypesConfig struct {
	PieceTypes []PieceTypeConfig `hcl:"piece_type,block"`
}

type PieceTypeConfig struct {
	Name string `hcl:"piece_name,label"`
}

type InitialStateConfig struct {
	Pieces []PiecesConfig `hcl:"pieces,block"`
}

type PiecesConfig struct {
	PlayerColor string         `hcl:"player_name,label"`
	Placements  hcl.Attributes `hcl:",remain"`
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

	err = mapstructure.Decode(funcs, &config.Functions)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) ToGame() (*game.GameState, error) {
	board, err := game.NewBoard(int(c.Board.Width), int(c.Board.Height))
	if err != nil {
		return nil, fmt.Errorf("creating board: %w", err)
	}

	pieceTypes := make(map[string]*game.PieceType, len(c.PieceTypes.PieceTypes))
	for _, pieceType := range c.PieceTypes.PieceTypes {
		pieceTypes[pieceType.Name] = &game.PieceType{
			Name: pieceType.Name,
		}
	}

	players := game.NewPlayers()

	state := &game.GameState{
		Board:   board,
		Players: players,
	}

	err = placePieces(state, c.InitialState.Pieces, pieceTypes)
	if err != nil {
		return nil, fmt.Errorf("placing initial pieces: %w", err)
	}

	return state, nil
}

func placePieces(state *game.GameState, pieces []PiecesConfig, pieceTypes map[string]*game.PieceType) error {
	for _, pieces := range pieces {
		color, err := game.ColorString(pieces.PlayerColor)
		if err != nil {
			return fmt.Errorf("parsing player color: %w", err)
		}
		player, err := state.GetPlayer(color)
		if err != nil {
			return fmt.Errorf("getting player: %w", err)
		}

		for _, piecePlacement := range pieces.Placements {
			squareString := piecePlacement.Name
			square, err := game.ParseSquare(squareString)
			if err != nil {
				return fmt.Errorf("parsing square: %w", err)
			}

			pieceTypeName, diags := piecePlacement.Expr.Value(nil)
			if diags.HasErrors() {
				return fmt.Errorf("parsing piece type: %w", diags)
			}
			pieceType := pieceTypes[pieceTypeName.AsString()]
			if pieceType == nil {
				return fmt.Errorf("piece type %q not defined", pieceTypeName)
			}

			piece := &game.Piece{
				Type:  pieceType,
				Owner: player,
			}

			err = state.Board.Place(piece, *square)
			if err != nil {
				return fmt.Errorf("placing a piece: %w", err)
			}
		}
	}
	return nil
}

func (c *Config) ToController() game.GameController {
	return c.Functions
}
