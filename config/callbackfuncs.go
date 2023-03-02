package config

import (
	"fmt"
	"image/color"

	"github.com/jostrzol/mess/game"
	"github.com/jostrzol/mess/game/piece"
	plr "github.com/jostrzol/mess/game/player"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

type CallbackFunctionsConfig struct {
	DecideWinnerFunc function.Function `mapstructure:"decide_winner"`
}

func (c CallbackFunctionsConfig) DecideWinner(state *game.State) (*plr.Player, error) {
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

func gameStateToCty(state *game.State) cty.Value {
	piecesPerPlayer := state.PiecesPerPlayer()
	players := make(map[string]cty.Value, len(state.Players))
	for _, player := range state.Players {
		pieces := piecesPerPlayer[player]
		players[color.Color.String()] = playerToCty(player, pieces)
	}
	return cty.ObjectVal(map[string]cty.Value{
		"players": cty.MapVal(players),
	})
}

func playerToCty(player *plr.Player, pieces []*piece.Piece) cty.Value {
	piecesCty := make([]cty.Value, len(pieces))
	for i, piece := range pieces {
		piecesCty[i] = pieceToCty(piece)
	}
	return cty.ObjectVal(map[string]cty.Value{
		"color":  cty.StringVal(color.Color.String()),
		"pieces": cty.ListVal(piecesCty),
	})
}

func pieceToCty(piece *piece.Piece) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"type":   cty.StringVal(piece.Type.Name),
		"square": cty.StringVal(piece.Square.String()),
	})
}

func playerFromCty(state *game.State, player cty.Value) (*plr.Player, error) {
	if player.IsNull() {
		return nil, nil
	}
	winnerColorStr := player.GetAttr("color").AsString()
	winnerColor, err := plr.ColorString(winnerColorStr)
	if err != nil {
		return nil, fmt.Errorf("parsing player color: %w", err)
	}
	winner, err := state.GetPlayer(winnerColor)
	if err != nil {
		return nil, fmt.Errorf("finding player in game: %w", err)
	}
	return winner, nil
}
