package mess_test

import (
	"testing"

	"github.com/jostrzol/mess/pkg/board/boardtest"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/stretchr/testify/assert"
)

func TestClone(t *testing.T) {
	board, err := mess.NewPieceBoard(2, 2)
	assert.NoError(t, err)

	piece := mess.NewPiece(King(t), nil)
	err = board.Place(piece, boardtest.NewSquare("A1"))
	assert.NoError(t, err)

	clone := board.Clone()

	err = piece.MoveTo(boardtest.NewSquare("A2"))
	assert.NoError(t, err)

	pieceA1, err := clone.At(boardtest.NewSquare("A1"))
	assert.NoError(t, err)
	assert.NotNil(t, pieceA1)
	pieceA2, err := clone.At(boardtest.NewSquare("A2"))
	assert.NoError(t, err)
	assert.Nil(t, pieceA2)
}
