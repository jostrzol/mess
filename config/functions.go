package config

import (
	"fmt"

	"github.com/jostrzol/mess/game"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

type FunctionsConfig struct {
	DecideWinnerFunc function.Function `mapstructure:"decide_winner"`
}

func (c FunctionsConfig) DecideWinner(state *game.GameState) (*game.Player, error) {
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
	players := make(map[string]cty.Value, len(state.Players))
	for _, player := range state.Players {
		players[player.Color.String()] = playerToCty(player)
	}
	return cty.ObjectVal(map[string]cty.Value{
		"players": cty.MapVal(players),
	})
}

func playerToCty(player *game.Player) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"color": cty.StringVal(player.Color.String()),
	})
}

func playerFromCty(state *game.GameState, player cty.Value) (*game.Player, error) {
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
