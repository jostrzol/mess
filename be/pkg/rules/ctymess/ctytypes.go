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
	"captures":          cty.Map(cty.Number),
	"forward_direction": Offset,
})

var Offset = cty.Tuple([]cty.Type{cty.Number, cty.Number})

var Piece = cty.Object(map[string]cty.Type{
	"type":   cty.String,
	"color":  cty.String,
	"square": cty.String,
})

var Record = cty.List(MoveGroup)

var Turn = cty.Object(map[string]cty.Type{
	"player": Player,
	"piece":  Piece,
	"src":    cty.String,
	"dst":    cty.String,
})

var MoveGroup = cty.Object(map[string]cty.Type{
	"name":   cty.String,
	"player": Player,
	"piece":  Piece,
	"src":    cty.String,
	"dst":    cty.String,
})

var Move = cty.Object(map[string]cty.Type{
	"name":    cty.String,
	"player":  Player,
	"piece":   Piece,
	"src":     cty.String,
	"dst":     cty.String,
	"options": cty.List(Option),
})

var SquareVec = cty.Object(map[string]cty.Type{
	"src": cty.String,
	"dst": cty.String,
})

var Option = cty.DynamicPseudoType
var Choice = cty.DynamicPseudoType

var Coords = Offset

var Board = cty.Object(map[string]cty.Type{
	"width":  cty.Number,
	"height": cty.Number,
})

var PieceType = cty.Object(map[string]cty.Type{
	"name": cty.String,
})
