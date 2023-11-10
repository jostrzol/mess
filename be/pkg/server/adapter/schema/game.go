package schema

import (
	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/jostrzol/mess/pkg/server/core/game"
)

type State struct {
	ID         uuid.UUID
	TurnNumber int
	Pieces     []Piece
	OptionTree *OptionNode
	IsMyTurn   bool
}

func StateFromDomain(s *game.State) *State {
	return &State{
		ID:         s.ID.UUID,
		TurnNumber: s.TurnNumber,
		Pieces:     piecesFromDomain(s.Board.AllPieces()),
		OptionTree: optionNodeFromDomain(s.OptionTree),
		IsMyTurn:   s.IsMyTurn,
	}
}

type Piece struct {
	Type   PieceType
	Color  string
	Square Square
}

func piecesFromDomain(pieces []*mess.Piece) []Piece {
	result := make([]Piece, 0, len(pieces))
	for _, piece := range pieces {
		result = append(result, Piece{
			Type:   pieceTypeFromDomain(piece.Type()),
			Color:  piece.Color().String(),
			Square: squareFromDomain(piece.Square()),
		})
	}
	return result
}

type PieceType struct {
	Name string
}

func pieceTypeFromDomain(pieceType *mess.PieceType) PieceType {
	return PieceType{
		Name: pieceType.Name(),
	}
}

type Square [2]int

func squareFromDomain(square board.Square) Square {
	x, y := square.ToCoords()
	return [2]int{x, y}
}

func (s Square) ToDomain() board.Square {
	return board.SquareFromCoords(s[0], s[1])
}

type SquareVec struct {
	From Square
	To   Square
}

func squareVecFromDomain(vec mess.SquareVec) SquareVec {
	return SquareVec{
		From: squareFromDomain(vec.From),
		To:   squareFromDomain(vec.To),
	}
}

func (s SquareVec) ToDomain() mess.SquareVec {
	return mess.SquareVec{
		From: s.From.ToDomain(),
		To:   s.To.ToDomain(),
	}
}
