package mess_test

import (
	"fmt"
	"testing"

	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/board/boardtest"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/event/eventtest"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/jostrzol/mess/pkg/mess/messtest"
	"github.com/stretchr/testify/suite"
)

func Rook(t *testing.T) *mess.PieceType {
	t.Helper()
	pieceType := mess.NewPieceType("rook")
	pieceType.AddMoveGenerator(mess.Motion{
		Name: "rook_motion",
		MoveGenerator: func(piece *mess.Piece) ([]board.Square, mess.MoveActionFunc) {
			result := make([]board.Square, 0)
			for _, offset := range []board.Offset{
				{X: 1, Y: 0},
				{X: -1, Y: 0},
				{X: 0, Y: 1},
				{X: 0, Y: -1},
			} {
				square := piece.Square().Offset(offset)
				for piece.Board().Contains(square) {
					result = append(result, square)
					square = square.Offset(offset)
				}
			}
			return result, nil
		},
	})
	return pieceType
}

func Knight(t *testing.T) *mess.PieceType {
	t.Helper()
	pieceType := mess.NewPieceType("knight")
	pieceType.AddMoveGenerator(messtest.OffsetMoveGenerator(t, []board.Offset{
		{X: 1, Y: 2},
		{X: 1, Y: -2},
		{X: -1, Y: 2},
		{X: -1, Y: -2},
		{X: 2, Y: 1},
		{X: 2, Y: -1},
		{X: -2, Y: 1},
		{X: -2, Y: -1},
	}...))
	return pieceType
}

func King(t *testing.T) *mess.PieceType {
	t.Helper()
	pieceType := mess.NewPieceType("king")
	pieceType.AddMoveGenerator(messtest.OffsetMoveGenerator(t, []board.Offset{
		{X: 1, Y: 0},
		{X: -1, Y: 0},
		{X: 0, Y: 1},
		{X: 0, Y: -1},
	}...))
	return pieceType
}

func Noones(pieceType *mess.PieceType) *mess.Piece {
	return mess.NewPiece(pieceType, nil)
}

type PieceSuite struct {
	suite.Suite
	board    *mess.PieceBoard
	observer *eventtest.MockObserver
}

func (s *PieceSuite) SetupTest() {
	board, err := mess.NewPieceBoard(4, 4)
	s.NoError(err)

	s.board = board
	s.observer = eventtest.NewMockObserver(s.T())
	s.board.Observe(s.observer)
}

func (s *PieceSuite) TestPlaceOn() {
	rook := Noones(Rook(s.T()))
	square := boardtest.NewSquare("B2")

	err := rook.PlaceOn(s.board, square)
	s.NoError(err)

	s.Equal(s.board, rook.Board())
	s.Equal(square, rook.Square())
	s.observer.ObservedMatch(mess.PiecePlaced{
		Piece:  rook,
		Board:  s.board,
		Square: square,
	})

	piece, err := s.board.At(square)
	s.NoError(err)
	s.Equal(rook, piece)
}

func (s *PieceSuite) TestPlaceOnReplace() {
	rook := Noones(Rook(s.T()))
	knight := Noones(Knight(s.T()))
	square := boardtest.NewSquare("B2")

	err := knight.PlaceOn(s.board, square)
	s.NoError(err)
	err = rook.PlaceOn(s.board, square)
	s.Error(err)

	s.Equal(s.board, knight.Board())
	s.Equal(square, knight.Square())
	s.False(rook.IsOnBoard())

	piece, err := s.board.At(square)
	s.NoError(err)

	s.Equal(knight, piece)
}

func (s *PieceSuite) TestMoves() {
	tests := []struct {
		pieceType     *mess.PieceType
		square        string
		expectedDests []string
	}{
		{
			pieceType:     Rook(s.T()),
			square:        "B2",
			expectedDests: []string{"B1", "B3", "B4", "A2", "C2", "D2"},
		},
		{
			pieceType:     Knight(s.T()),
			square:        "B2",
			expectedDests: []string{"A4", "C4", "D1", "D3"},
		},
	}
	for _, tt := range tests {
		s.SetupTest()
		s.Run(fmt.Sprintf("%v at %v", tt.pieceType, tt.square), func() {
			piece := Noones(tt.pieceType)
			piece.PlaceOn(s.board, boardtest.NewSquare(tt.square))

			moves := piece.Moves()

			messtest.MovesMatch(s.T(), moves, messtest.MovesMatcher(piece, tt.expectedDests...))
		})
	}
}

func (s *PieceSuite) TestMove() {
	startSquare := boardtest.NewSquare("B2")
	endSquare := boardtest.NewSquare("C4")

	knight := Noones(Knight(s.T()))
	knight.PlaceOn(s.board, startSquare)

	s.observer.Reset()
	err := knight.MoveTo(endSquare)
	s.NoError(err)

	s.Equal(endSquare, knight.Square())
	s.observer.ObservedMatch(mess.PieceMoved{
		Piece: knight,
		From:  startSquare,
		To:    endSquare,
	})

	empty, err := s.board.At(startSquare)
	s.NoError(err)
	s.Nil(empty)

	piece, err := s.board.At(endSquare)
	s.NoError(err)
	s.Equal(knight, piece)
}

func (s *PieceSuite) TestMoveReplace() {
	startSquare := boardtest.NewSquare("B2")
	endSquare := boardtest.NewSquare("C4")

	players := mess.NewPlayers(s.board)

	knight := mess.NewPiece(Knight(s.T()), players[color.White])
	knight.PlaceOn(s.board, startSquare)
	rook := mess.NewPiece(Rook(s.T()), players[color.Black])
	rook.PlaceOn(s.board, endSquare)

	s.observer.Reset()
	err := knight.MoveTo(endSquare)
	s.NoError(err)

	s.Equal(endSquare, knight.Square())
	s.False(rook.IsOnBoard())
	s.observer.ObservedMatch(mess.PieceMoved{
		Piece: knight,
		From:  startSquare,
		To:    endSquare,
	}, mess.PieceCaptured{
		Piece:        rook,
		CapturedBy:   players[color.White],
		CapturedFrom: players[color.Black],
	}, mess.PieceRemoved{
		Piece:  rook,
		Square: endSquare,
	})

	empty, err := s.board.At(startSquare)
	s.NoError(err)
	s.Nil(empty)

	piece, err := s.board.At(endSquare)
	s.NoError(err)
	s.Equal(knight, piece)
}

func (s *PieceSuite) TestMoveOutOfBounds() {
	startSquare := boardtest.NewSquare("B2")
	endSquare := boardtest.NewSquare("Z1")

	knight := Noones(Knight(s.T()))
	knight.PlaceOn(s.board, startSquare)

	err := knight.MoveTo(endSquare)
	s.Error(err)
}

func (s *PieceSuite) TestIsOnBoard() {
	knight := Noones(Knight(s.T()))
	knight.PlaceOn(s.board, boardtest.NewSquare("A1"))

	s.True(knight.IsOnBoard())
}

func (s *PieceSuite) TestIsNotOnBoard() {
	knight := Noones(Knight(s.T()))
	s.False(knight.IsOnBoard())
}

func TestPieceSuite(t *testing.T) {
	suite.Run(t, new(PieceSuite))
}
