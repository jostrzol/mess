package ctyconv

import (
	"github.com/jostrzol/mess/game"
	"github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/piece"
	plr "github.com/jostrzol/mess/game/player"
	"github.com/zclconf/go-cty/cty"
)

func GameStateToCty(state *game.State) cty.Value {
	piecesPerPlayer := state.PiecesPerPlayer()
	players := make(map[string]cty.Value, len(state.Players))
	for _, player := range state.Players {
		pieces := piecesPerPlayer[player]
		players[player.Color().String()] = PlayerToCty(player, pieces)
	}
	return cty.ObjectVal(map[string]cty.Value{
		"players": cty.MapVal(players),
	})
}

func PlayerToCty(player *plr.Player, pieces []*piece.Piece) cty.Value {
	piecesCty := make([]cty.Value, len(pieces))
	for i, piece := range pieces {
		piecesCty[i] = PieceToCty(piece)
	}
	return cty.ObjectVal(map[string]cty.Value{
		"color":  cty.StringVal(player.Color().String()),
		"pieces": cty.ListVal(piecesCty),
	})
}

func PieceToCty(piece *piece.Piece) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"type":   cty.StringVal(piece.Type.Name),
		"square": cty.StringVal(piece.Square.String()),
	})
}

func SquareToCty(square *board.Square) cty.Value {
	return cty.StringVal(square.String())
}
