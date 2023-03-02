package piece_test

import (
	"testing"

	brd "github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/board/boardtest"
	"github.com/jostrzol/mess/game/piece"
	"github.com/jostrzol/mess/game/piece/piecetest"
	"github.com/stretchr/testify/assert"
)

func TestGenerateMotionsRook(t *testing.T) {
	board, err := brd.NewBoard[*piece.Piece](4, 4)
	assert.NoError(t, err)

	rook := piecetest.Noones(piecetest.Rook(t))
	rook.PlaceOn(board, boardtest.NewSquare("B2"))

	motions := rook.GenerateMotions()
	assert.ElementsMatch(t, motions, []brd.Square{
		boardtest.NewSquare("B1"),
		boardtest.NewSquare("B3"),
		boardtest.NewSquare("B4"),
		boardtest.NewSquare("A2"),
		boardtest.NewSquare("C2"),
		boardtest.NewSquare("D2"),
	})
}

func TestGenerateMotionsKnight(t *testing.T) {
	board, err := brd.NewBoard[*piece.Piece](4, 4)
	assert.NoError(t, err)

	knight := piecetest.Noones(piecetest.Knight(t))
	knight.PlaceOn(board, boardtest.NewSquare("B2"))

	motions := knight.GenerateMotions()
	assert.ElementsMatch(t, motions, []brd.Square{
		boardtest.NewSquare("A4"),
		boardtest.NewSquare("C4"),
		boardtest.NewSquare("D1"),
		boardtest.NewSquare("D3"),
	})
}
