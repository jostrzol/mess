package mess

import "github.com/jostrzol/mess/pkg/board"

type Choice interface{ GenerateOptions() []Option }
type Option interface{ String() string }

// Piece type choice

type PieceTypeChoice struct {
	PieceTypes []*PieceType
}

func (ptc PieceTypeChoice) GenerateOptions() []Option {
	result := make([]Option, len(ptc.PieceTypes))
	for i, pieceType := range ptc.PieceTypes {
		result[i] = PieceTypeOption{pieceType}
	}
	return result
}

type PieceTypeOption struct{ *PieceType }

func (pto PieceTypeOption) String() string {
	return pto.name
}

// Square choice

type SquareChoice struct {
	Squares []board.Square
}

func (sc SquareChoice) GenerateOptions() []Option {
	result := make([]Option, len(sc.Squares))
	for i, square := range sc.Squares {
		result[i] = SquareOption{square}
	}
	return result
}

type SquareOption struct{ board.Square }

func (so SquareOption) String() string {
	return so.String()
}
