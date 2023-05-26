package integration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jostrzol/mess/config"
	"github.com/jostrzol/mess/config/configtest"
	"github.com/jostrzol/mess/pkg/board/boardtest"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/stretchr/testify/assert"
)

func place(t *testing.T, game *mess.Game, color color.Color, pieceName string, square string) {
	t.Helper()
	pieceType, err := game.GetPieceType(pieceName)
	assert.NoError(t, err)
	owner := game.Player(color)
	piece := mess.NewPiece(pieceType, owner)
	game.Board().Place(piece, boardtest.NewSquare(square))
}

func TestMoves(t *testing.T) {
	type piece struct {
		color  color.Color
		name   string
		square string
	}
	type move struct {
		from string
		to   string
	}
	type test struct {
		name        string
		initState   []piece
		whenMoves   []move
		expectMoves []move
	}

	tests := []test{
		{
			name: "castling",
			initState: []piece{
				{color.White, "king", "E1"},
				{color.White, "rook", "A1"},
				{color.White, "rook", "H1"},
			},
			expectMoves: []move{
				{"E1", "C1"},
				{"E1", "G1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game, err := config.DecodeConfig(RulesFile, configtest.RandomInteractor{}, false)
			assert.NoError(t, err)

			for _, piece := range tt.initState {
				place(t, game, piece.color, piece.name, piece.square)
			}

			for _, move := range tt.whenMoves {
				moveMade := false
				validMoves := game.ValidMoves()
				for _, validMove := range validMoves {
					if validMove.From.String() == move.from && validMove.To.String() == move.to {
						validMove.Perform()
						moveMade = true
						break
					}
				}
				if !moveMade {
					t.Errorf("precondition move not valid: %v -> %v", move.from, move.to)
				}
			}

			validMoves := game.ValidMoves()
			notFound := make([]move, 0)
			for _, move := range tt.expectMoves {
				found := false
				for _, validMove := range validMoves {
					if validMove.From.String() == move.from && validMove.To.String() == move.to {
						found = true
						break
					}
				}
				if !found {
					notFound = append(notFound, move)
				}
			}

			if len(notFound) > 0 {
				var msg strings.Builder
				msg.WriteString("expected moves not found: [")
				for i, move := range notFound {
					if i != 0 {
						msg.WriteString(", ")
					}
					msg.WriteString(fmt.Sprintf("%s -> %s", move.from, move.to))
				}
				msg.WriteString("]")
				t.Fatal(msg.String())
			}
		})
	}
}
