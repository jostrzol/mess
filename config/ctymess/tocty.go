package ctymess

import (
	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/zclconf/go-cty/cty"
)

func GameStateToCty(state *mess.State) cty.Value {
	players := make(map[string]cty.Value, len(state.Players()))
	for player := range state.Players() {
		players[player.Color().String()] = PlayerToCty(player)
	}
	return cty.ObjectVal(map[string]cty.Value{
		"players":        cty.MapVal(players),
		"current_player": PlayerToCty(state.CurrentPlayer()),
	})
}

func PlayerToCty(player *mess.Player) cty.Value {
	piecesCty := make([]cty.Value, 0, len(player.Pieces()))
	for piece := range player.Pieces() {
		piecesCty = append(piecesCty, PieceToCty(piece))
	}
	return cty.ObjectVal(map[string]cty.Value{
		"color":             cty.StringVal(player.Color().String()),
		"pieces":            cty.ListVal(piecesCty),
		"forward_direction": OffsetToCty(player.ForwardDirection()),
	})
}

func PieceToCty(piece *mess.Piece) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"type":   cty.StringVal(piece.Type().Name()),
		"color":  cty.StringVal(piece.Color().String()),
		"square": cty.StringVal(piece.Square().String()),
	})
}

func SquareToCty(square board.Square) cty.Value {
	return cty.StringVal(square.String())
}

func OffsetToCty(offset board.Offset) cty.Value {
	return cty.TupleVal([]cty.Value{
		cty.NumberIntVal(int64(offset.X)),
		cty.NumberIntVal(int64(offset.Y)),
	})
}
