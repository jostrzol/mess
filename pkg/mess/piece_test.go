package mess

import (
	"testing"

	"github.com/jostrzol/mess/pkg/board"
	brd "github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/board/boardtest"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/event/eventtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func Rook(t *testing.T) *PieceType {
	t.Helper()
	pieceType := NewPieceType("rook")
	pieceType.AddMoveGenerator(func(piece *Piece) []board.Square {
		result := make([]board.Square, 0)
		for _, offset := range []board.Offset{
			{X: 1, Y: 0},
			{X: -1, Y: 0},
			{X: 0, Y: 1},
			{X: 0, Y: -1},
		} {
			square := piece.Square().Offset(offset)
			for piece.Board().Contains(square) {
				result = append(result, *square)
				square = square.Offset(offset)
			}
		}
		return result
	})
	return pieceType
}

func Knight(t *testing.T) *PieceType {
	t.Helper()
	pieceType := NewPieceType("knight")
	pieceType.AddMoveGenerator(offsetMoveGenerator(t, []board.Offset{
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

func King(t *testing.T) *PieceType {
	t.Helper()
	pieceType := NewPieceType("king")
	pieceType.AddMoveGenerator(offsetMoveGenerator(t, []board.Offset{
		{X: 1, Y: 0},
		{X: -1, Y: 0},
		{X: 0, Y: 1},
		{X: 0, Y: -1},
	}...))
	return pieceType
}

func Noones(pieceType *PieceType) *Piece {
	return NewPiece(pieceType, nil)
}

type PieceSuite struct {
	suite.Suite
	board    *PieceBoard
	observer *eventtest.MockObserver
}

func (s *PieceSuite) SetupTest() {
	board, err := NewPieceBoard(4, 4)
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
	s.Equal(*square, *rook.Square())
	s.observer.ObservedMatch(PiecePlaced{
		Piece:  rook,
		Board:  s.board,
		Square: *square,
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
	s.Equal(*square, *knight.Square())
	s.False(rook.IsOnBoard())

	piece, err := s.board.At(square)
	s.NoError(err)

	s.Equal(knight, piece)
}

func (s *PieceSuite) TestGenerateMovesRook() {
	rook := Noones(Rook(s.T()))
	rook.PlaceOn(s.board, boardtest.NewSquare("B2"))

	moves := rook.ValidMoves()
	s.ElementsMatch(moves, movesFromDests(rook,
		*boardtest.NewSquare("B1"),
		*boardtest.NewSquare("B3"),
		*boardtest.NewSquare("B4"),
		*boardtest.NewSquare("A2"),
		*boardtest.NewSquare("C2"),
		*boardtest.NewSquare("D2"),
	))
}

func (s *PieceSuite) TestGenerateMovesKnight() {
	knight := Noones(Knight(s.T()))
	knight.PlaceOn(s.board, boardtest.NewSquare("B2"))

	moves := knight.ValidMoves()
	s.ElementsMatch(moves, movesFromDests(knight,
		*boardtest.NewSquare("A4"),
		*boardtest.NewSquare("C4"),
		*boardtest.NewSquare("D1"),
		*boardtest.NewSquare("D3"),
	))
}

func movesFromDests(piece *Piece, destinations ...board.Square) []Move {
	result := make([]Move, len(destinations))
	for i, destination := range destinations {
		result[i].Piece = piece
		result[i].From = *piece.Square()
		result[i].To = destination
	}
	return result
}

func (s *PieceSuite) TestMove() {
	startSquare := boardtest.NewSquare("B2")
	endSquare := boardtest.NewSquare("C4")

	knight := Noones(Knight(s.T()))
	knight.PlaceOn(s.board, startSquare)

	s.observer.Reset()
	err := knight.MoveTo(endSquare)
	s.NoError(err)

	s.Equal(*endSquare, *knight.Square())
	s.observer.ObservedMatch(PieceMoved{
		Piece: knight,
		From:  *startSquare,
		To:    *endSquare,
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

	players := NewPlayers(s.board)

	knight := NewPiece(Knight(s.T()), players[color.White])
	knight.PlaceOn(s.board, startSquare)
	rook := NewPiece(Rook(s.T()), players[color.Black])
	rook.PlaceOn(s.board, endSquare)

	s.observer.Reset()
	err := knight.MoveTo(endSquare)
	s.NoError(err)

	s.Equal(*endSquare, *knight.Square())
	s.False(rook.IsOnBoard())
	s.observer.ObservedMatch(PieceMoved{
		Piece: knight,
		From:  *startSquare,
		To:    *endSquare,
	}, PieceCaptured{
		Piece:        rook,
		CapturedBy:   players[color.White],
		CapturedFrom: players[color.Black],
	}, PieceRemoved{
		Piece:  rook,
		Square: *endSquare,
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

func staticMoveGenerator(t *testing.T, strings ...string) MoveGenerator {
	t.Helper()
	return func(piece *Piece) []brd.Square {
		destinations := make([]brd.Square, 0, len(strings))
		for _, squareStr := range strings {
			square, err := brd.NewSquare(squareStr)
			assert.NoError(t, err)
			destinations = append(destinations, *square)
		}
		return destinations
	}
}

func offsetMoveGenerator(t *testing.T, offsets ...board.Offset) MoveGenerator {
	t.Helper()
	return func(piece *Piece) []brd.Square {
		destinations := make([]brd.Square, 0, len(offsets))
		for _, offset := range offsets {
			square := piece.Square().Offset(offset)
			if piece.Board().Contains(square) {
				destinations = append(destinations, *square)
			}
		}
		return destinations
	}
}

func assertSquaresMatch(t *testing.T, actual []brd.Square, expected ...string) {
	assert.Len(t, actual, len(expected))
	for _, str := range expected {
		square, err := brd.NewSquare(str)
		assert.NoError(t, err)
		assert.Containsf(t, actual, *square, "%v doesnt contain square %v", actual, square)
	}
}

func TestGenerateMovesZero(t *testing.T) {
	generators := moveGenerators([]MoveGenerator{})
	moves := generators.Generate(nil)
	assert.Empty(t, moves)
}

func TestGenerateMovesOne(t *testing.T) {
	generators := moveGenerators([]MoveGenerator{
		staticMoveGenerator(t, "A1"),
	})
	moves := generators.Generate(nil)
	assertSquaresMatch(t, moves, "A1")
}

func TestGenerateMovesTwo(t *testing.T) {
	generators := moveGenerators([]MoveGenerator{
		staticMoveGenerator(t, "A1"),
		staticMoveGenerator(t, "B1"),
	})
	moves := generators.Generate(nil)
	assertSquaresMatch(t, moves, "A1", "B1")
}

func TestGenerateMovesTwoOverlapping(t *testing.T) {
	generators := moveGenerators([]MoveGenerator{
		staticMoveGenerator(t, "A1"),
		staticMoveGenerator(t, "A1"),
	})
	moves := generators.Generate(nil)
	assertSquaresMatch(t, moves, "A1")
}

func TestGenerateMovesMany(t *testing.T) {
	generators := moveGenerators([]MoveGenerator{
		staticMoveGenerator(t, "A1", "B2"),
		staticMoveGenerator(t, "C5"),
		staticMoveGenerator(t, "B2", "D4", "C5"),
	})
	moves := generators.Generate(nil)
	assertSquaresMatch(t, moves, "A1", "B2", "C5", "D4")
}
