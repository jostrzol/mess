package schema

import (
	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/jostrzol/mess/pkg/server/core/room"
)

type State struct {
	Pieces     []Piece
	OptionTree interface{}
}

func StateFromDomain(s *room.State) *State {
	return &State{
		Pieces:     piecesFromDomain(s.Board.AllPieces()),
		OptionTree: optionNodeFromDomain(s.OptionTree),
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

type MoveGroup struct {
	From Square
	To   Square
}

func moveGroupFromDomain(moveGroup *mess.MoveGroup) MoveGroup {
	return MoveGroup{
		From: squareFromDomain(moveGroup.From),
		To:   squareFromDomain(moveGroup.To),
	}
}
