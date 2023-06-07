package ctymess

import (
	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/zclconf/go-cty/cty"
)

func StateToCty(state *mess.State) cty.Value {
	players := make(map[string]cty.Value, len(state.Players()))
	for player := range state.Players() {
		players[player.Color().String()] = PlayerToCty(player)
	}
	return cty.ObjectVal(map[string]cty.Value{
		"players":        cty.MapVal(players),
		"current_player": players[state.CurrentPlayer().Color().String()],
		"record":         RecordToCty(state.Record()),
	})
}

func PlayerToCty(player *mess.Player) cty.Value {
	pieces := make([]cty.Value, 0, len(player.Pieces()))
	for piece := range player.Pieces() {
		pieces = append(pieces, PieceToCty(piece))
	}
	return cty.ObjectVal(map[string]cty.Value{
		"color":             cty.StringVal(player.Color().String()),
		"pieces":            listOrEmpty(Piece, pieces),
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

func RecordToCty(record []mess.RecordedMove) cty.Value {
	result := make([]cty.Value, 0, len(record))
	for _, move := range record {
		captures := make([]cty.Value, 0, len(move.Captures))
		for piece := range move.Captures {
			captures = append(captures, PieceToCty(piece))
		}
		moveCty := cty.ObjectVal(map[string]cty.Value{
			"piece":    PieceToCty(move.Piece),
			"player":   PlayerToCty(move.Piece.Owner()),
			"src":      cty.StringVal(move.From.String()),
			"dst":      cty.StringVal(move.To.String()),
			"captures": listOrEmpty(Piece, captures),
		})
		result = append(result, moveCty)
	}
	return listOrEmpty(RecordedMove, result)
}

func MoveToCty(move *mess.Move) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"name":   cty.StringVal(move.Name),
		"piece":  PieceToCty(move.Piece),
		"player": PlayerToCty(move.Piece.Owner()),
		"src":    cty.StringVal(move.From.String()),
		"dst":    cty.StringVal(move.To.String()),
	})
}

func listOrEmpty(listType cty.Type, slice []cty.Value) cty.Value {
	if len(slice) == 0 {
		return cty.ListValEmpty(listType)
	}
	return cty.ListVal(slice)
}

func BoardToCty(board *mess.PieceBoard) cty.Value {
	width, height := board.Size()
	return cty.ObjectVal(map[string]cty.Value{
		"width":  cty.NumberIntVal(int64(width)),
		"height": cty.NumberIntVal(int64(height)),
	})
}

func PieceTypesToCty(pieceTypes []*mess.PieceType) cty.Value {
	ctyPieceTypes := make([]cty.Value, 0, len(pieceTypes))
	for _, pt := range pieceTypes {
		ctyPt := cty.ObjectVal(map[string]cty.Value{
			"name": cty.StringVal(pt.Name()),
		})
		ctyPieceTypes = append(ctyPieceTypes, ctyPt)
	}
	return cty.ListVal(ctyPieceTypes)
}
