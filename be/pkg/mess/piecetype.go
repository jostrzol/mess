package mess

import (
	"unicode"
	"unicode/utf8"

	"github.com/jostrzol/mess/pkg/color"
)

type PieceType struct {
	name         string
	presentation map[color.Color]Presentation
	motions      chainMotions
}

type Presentation struct {
	Symbol rune
	Icon   AssetKey
	Rotate bool
}

func NewPieceType(name string) *PieceType {
	return &PieceType{
		name: name,
		presentation: map[color.Color]Presentation{
			color.Black: {Symbol: defaultSymbol(color.Black, name)},
			color.White: {Symbol: defaultSymbol(color.White, name)},
		},
		motions: make(chainMotions, 0),
	}
}

func (t *PieceType) Name() string {
	return t.name
}

func (t *PieceType) SetPresentation(color color.Color, presentation Presentation) {
	if presentation.Symbol == 0 {
		presentation.Symbol = defaultSymbol(color, t.Name())
	}
	t.presentation[color] = presentation
}

func (t *PieceType) Presentation(color color.Color) Presentation {
	return t.presentation[color]
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
