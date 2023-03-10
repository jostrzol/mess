package ctymess

import "github.com/zclconf/go-cty/cty"

var Game = cty.Object(map[string]cty.Type{
	"players":        cty.Map(Player),
	"current_player": Player,
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
