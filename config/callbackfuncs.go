package config

import (
	"fmt"

	"github.com/jostrzol/mess/game"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

type CallbackFunctionsConfig struct {
	DecideWinnerFunc function.Function `mapstructure:"decide_winner"`
}

func (c CallbackFunctionsConfig) DecideWinner(state *game.GameState) (*game.Player, error) {
	ctyState := gameStateToCty(state)
	ctyWinner, err := c.DecideWinnerFunc.Call([]cty.Value{ctyState})
	if err != nil {
		return nil, fmt.Errorf("calling user-defined function: %w", err)
	}
	winner, err := playerFromCty(state, ctyWinner)
	if err != nil {
		return nil, fmt.Errorf("getting winner: %w", err)
	}
	return winner, nil
}

func gameStateToCty(state *game.GameState) cty.Value {
	piecesPerPlayer := state.PiecesPerPlayer()
	players := make(map[string]cty.Value, len(state.Players))
	for _, player := range state.Players {
		pieces := piecesPerPlayer[player]
		players[player.Color.String()] = playerToCty(player, pieces)
	}
	return cty.ObjectVal(map[string]cty.Value{
		"players": cty.MapVal(players),
	})
}

func playerToCty(player *game.Player, pieces []*game.PieceOnSquare) cty.Value {
	piecesCty := make([]cty.Value, len(pieces))
	for i, piece := range pieces {
		piecesCty[i] = pieceOnSquareToCty(piece)
	}
	return cty.ObjectVal(map[string]cty.Value{
		"color":  cty.StringVal(player.Color.String()),
		"pieces": cty.ListVal(piecesCty),
	})
}

func pieceOnSquareToCty(pieceOnSquare *game.PieceOnSquare) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"type":   cty.StringVal(pieceOnSquare.Piece.Type.Name),
		"square": cty.StringVal(pieceOnSquare.Square.String()),
	})
}

func playerFromCty(state *game.GameState, player cty.Value) (*game.Player, error) {
	if player.IsNull() {
		return nil, nil
	}
	winnerColorStr := player.GetAttr("color").AsString()
	winnerColor, err := game.ColorString(winnerColorStr)
	if err != nil {
		return nil, fmt.Errorf("parsing player color: %w", err)
	}
	winner, err := state.GetPlayer(winnerColor)
	if err != nil {
		return nil, fmt.Errorf("finding player in game: %w", err)
	}
	return winner, nil
}
