package mess

import (
	"unicode"
	"unicode/utf8"

	"github.com/jostrzol/mess/pkg/color"
)

type PieceType struct {
	name           string
	representation map[color.Color]Representation
	motions        chainMotions
}

type Representation struct {
	Symbol rune
	Icon   AssetKey
}

func NewPieceType(name string) *PieceType {
	return &PieceType{
		name: name,
		representation: map[color.Color]Representation{
			color.Black: {Symbol: defaultSymbol(color.Black, name)},
			color.White: {Symbol: defaultSymbol(color.White, name)},
		},
		motions: make(chainMotions, 0),
	}
}

func (t *PieceType) Name() string {
	return t.name
}

func (t *PieceType) SetRepresentation(color color.Color, representation Representation) {
	if representation.Symbol == 0 {
		representation.Symbol = defaultSymbol(color, t.Name())
	}
	t.representation[color] = representation
}

func (t *PieceType) Representation(color color.Color) Representation {
	return t.representation[color]
}

func defaultSymbol(col color.Color, name string) rune {
	r, _ := utf8.DecodeRuneInString(name)
	if r == utf8.RuneError {
		return rune('?')
	}
	switch col {
	case color.Black:
		return unicode.ToLower(r)
	case color.White:
		return unicode.ToUpper(r)
	default:
		panic("unreachable")
	}
}

func (t *PieceType) String() string {
	return t.Name()
}

func (t *PieceType) AddMotion(motion Motion) {
	t.motions = append(t.motions, motion)
}

func (t *PieceType) moves(piece *Piece) []*MoveGroup {
	return t.motions.Generate(piece)
}
