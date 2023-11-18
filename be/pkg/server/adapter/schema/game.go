package schema

import (
	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/jostrzol/mess/pkg/server/core/game"
	"github.com/jostrzol/mess/pkg/server/core/id"
)

type StaticData struct {
	ID        uuid.UUID
	BoardSize BoardSize
	MyColor   string
}

type BoardSize struct {
	Width  int
	Height int
}

func StaticDataFromDomain(s *game.StaticData) *StaticData {
	return &StaticData{
		ID:        s.ID.UUID,
		BoardSize: BoardSize(s.BoardSize),
		MyColor:   s.MyColor.String(),
	}
}

type Resolution struct {
	Status string
}

func ResolutionFromDomain(session id.Session, r *game.Resolution) *Resolution {
	var status string
	switch {
	case !r.IsResolved:
		status = "Unresolved"
	case r.Winner == session:
		status = "Win"
	case r.Winner.IsZero():
		status = "Draw"
	default:
		status = "Defeat"
	}
	return &Resolution{Status: status}
}

type State struct {
	ID         uuid.UUID
	TurnNumber int
	Pieces     []Piece
	IsMyTurn   bool
}

func StateFromDomain(session id.Session, s *game.State) *State {
	return &State{
		ID:         s.ID.UUID,
		TurnNumber: s.TurnNumber,
		Pieces:     piecesFromDomain(s.Board.AllPieces()),
		IsMyTurn:   s.CurrentPlayer == session,
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
	Name           string
	Representation map[string]Representation
}

func pieceTypeFromDomain(pieceType *mess.PieceType) PieceType {
	return PieceType{
		Name: pieceType.Name(),
		Representation: map[string]Representation{
			color.Black.String(): representationFromDomain(pieceType.Representation(color.Black)),
			color.White.String(): representationFromDomain(pieceType.Representation(color.White)),
		},
	}
}

type Representation struct {
	Symbol string
	Icon   string `json:",omitempty"`
}

func representationFromDomain(representation mess.Representation) Representation {
	return Representation{
		Symbol: string(representation.Symbol),
		Icon:   string(representation.Icon),
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
