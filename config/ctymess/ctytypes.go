package ctymess

import "github.com/zclconf/go-cty/cty"

var Game = cty.Object(map[string]cty.Type{
	"players":        cty.Map(Player),
	"current_player": Player,
	"record":         Record,
})

var Player = cty.Object(map[string]cty.Type{
	"color":             cty.String,
	"pieces":            cty.List(Piece),
	"forward_direction": Offset,
})

var Offset = cty.Tuple([]cty.Type{cty.Number, cty.Number})

var Piece = cty.Object(map[string]cty.Type{
	"type":   cty.String,
	"color":  cty.String,
	"square": cty.String,
})

var Record = cty.List(Move)

var RecordedMove = cty.Object(map[string]cty.Type{
	"player":   Player,
	"piece":    Piece,
	"src":      cty.String,
	"dst":      cty.String,
	"captures": cty.List(Piece),
})

var Move = cty.Object(map[string]cty.Type{
	"player": Player,
	"piece":  Piece,
	"src":    cty.String,
	"dst":    cty.String,
})

var Coords = Offset

var Board = cty.Object(map[string]cty.Type{
	"width":  cty.Number,
	"height": cty.Number,
})
