package ctymess

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/zclconf/go-cty/cty"
)

func StateToCty(state *mess.State) cty.Value {
	players := make(map[string]cty.Value, len(state.Players()))
	for _, player := range state.Players() {
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
	for _, piece := range player.Pieces() {
		pieces = append(pieces, PieceToCty(piece))
	}
	captures := player.CapturesCountByType()
	capturesCty := make(map[string]cty.Value, len(captures))
	for pieceType, count := range captures {
		capturesCty[pieceType.Name()] = cty.NumberIntVal(int64(count))
	}
	return cty.ObjectVal(map[string]cty.Value{
		"color":             cty.StringVal(player.Color().String()),
		"pieces":            listOrEmpty(Piece, pieces),
		"captures":          mapOrEmpty(cty.Number, capturesCty),
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

func RecordToCty(record []mess.Turn) cty.Value {
	result := make([]cty.Value, 0, len(record))
	for _, turn := range record {
		move := turn.FirstMove()
		if move != nil {
			moveCty := cty.ObjectVal(map[string]cty.Value{
				"piece":  PieceToCty(move.Piece),
				"player": PlayerToCty(move.Piece.Owner()),
				"src":    cty.StringVal(move.From.String()),
				"dst":    cty.StringVal(move.To.String()),
			})
			result = append(result, moveCty)
		} else {
			// TODO: append something more useful
			result = append(result, cty.DynamicVal)
		}
	}
	return listOrEmpty(Turn, result)
}

func MoveToCty(move *mess.Move) cty.Value {
	if move == nil {
		return cty.NullVal(Move)
	}
	result := cty.ObjectVal(map[string]cty.Value{
		"name":    cty.StringVal(move.Name),
		"piece":   PieceToCty(move.Piece),
		"player":  PlayerToCty(move.Piece.Owner()),
		"src":     cty.StringVal(move.From.String()),
		"dst":     cty.StringVal(move.To.String()),
		"options": OptionsToCty(move.Options),
	})
	return result
}

func MoveGroupToCty(moveGroup *mess.MoveGroup) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"name":   cty.StringVal(moveGroup.Name),
		"piece":  PieceToCty(moveGroup.Piece),
		"player": PlayerToCty(moveGroup.Piece.Owner()),
		"src":    cty.StringVal(moveGroup.From.String()),
		"dst":    cty.StringVal(moveGroup.To.String()),
	})
}

func SquareVecToCty(vec mess.SquareVec) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"src": cty.StringVal(vec.From.String()),
		"dst": cty.StringVal(vec.To.String()),
	})
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
		ctyPt := PieceTypeToCty(pt)
		ctyPieceTypes = append(ctyPieceTypes, ctyPt)
	}
	return cty.ListVal(ctyPieceTypes)
}

func PieceTypeToCty(pieceType *mess.PieceType) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"name": cty.StringVal(pieceType.Name()),
	})
}

func OptionsToCty(options []mess.Option) cty.Value {
	if options == nil {
		return cty.NullVal(cty.Tuple([]cty.Type{Option}))
	}

	result := make([]cty.Value, 0, len(options))
	for _, option := range options {
		result = append(result, OptionToCty(option))
	}
	return tupleOrEmpty(result)
}

func OptionToCty(option mess.Option) cty.Value {
	switch opt := option.(type) {
	case mess.PieceTypeOption:
		return cty.ObjectVal(map[string]cty.Value{
			"type":       cty.StringVal("piece_type"),
			"piece_type": PieceTypeToCty(opt.PieceType),
		})
	case mess.SquareOption:
		return cty.ObjectVal(map[string]cty.Value{
			"type":   cty.StringVal("square"),
			"square": SquareToCty(opt.Square),
		})
	case mess.MoveOption:
		return cty.ObjectVal(map[string]cty.Value{
			"type": cty.StringVal("move"),
			"move": SquareVecToCty(opt.SquareVec),
		})
	case mess.UnitOption:
		return cty.ObjectVal(map[string]cty.Value{
			"type": cty.StringVal("unit"),
		})
	default:
		err := fmt.Errorf("invalid option type %T", option)
		panic(err)
	}
}

func mapOrEmpty(mapType cty.Type, mapValue map[string]cty.Value) cty.Value {
	if len(mapValue) == 0 {
		return cty.MapValEmpty(mapType)
	}
	return cty.MapVal(mapValue)
}

func listOrEmpty(listType cty.Type, slice []cty.Value) cty.Value {
	if len(slice) == 0 {
		return cty.ListValEmpty(listType)
	}
	return cty.ListVal(slice)
}

func tupleOrEmpty(slice []cty.Value) cty.Value {
	if len(slice) == 0 {
		return cty.EmptyTupleVal
	}
	return cty.TupleVal(slice)
}
