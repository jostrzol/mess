package game

import "fmt"

type GameState struct {
	Board   Board
	Players map[Color]*Player
}

func (g *GameState) GetPlayer(color Color) (*Player, error) {
	player, ok := g.Players[color]
	if !ok {
		return nil, fmt.Errorf("player of color %q not found", color)
	}
	return player, nil
}

type GameController interface {
	DecideWinner(state *GameState) (*Player, error)
}
