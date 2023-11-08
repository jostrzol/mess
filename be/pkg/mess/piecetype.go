package mess

import (
	"unicode"
	"unicode/utf8"
)

type PieceType struct {
	name    string
	symbols symbols
	motions chainMotions
}

type symbols struct {
	white rune
	black rune
}

func NewPieceType(name string) *PieceType {
	return &PieceType{
		name:    name,
		motions: make(chainMotions, 0),
	}
}

func (t *PieceType) Name() string {
	return t.name
}

func (t *PieceType) SymbolWhite() rune {
	if t.symbols.white == 0 {
		return unicode.ToUpper(t.firstNameRune())
	}
	return t.symbols.white
}

func (t *PieceType) SymbolBlack() rune {
	if t.symbols.black == 0 {
		return unicode.ToLower(t.firstNameRune())
	}
	return t.symbols.black
}

func (t *PieceType) firstNameRune() rune {
	r, _ := utf8.DecodeRuneInString(t.name)
	if r == utf8.RuneError {
		return rune('?')
	}
	return r
}

func (t *PieceType) SetSymbols(symbolWhite rune, symbolBlack rune) {
	t.symbols = symbols{
		white: symbolWhite,
		black: symbolBlack,
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
