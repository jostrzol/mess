package integration

import (
	"testing"

	"github.com/jostrzol/mess/config"
	"github.com/jostrzol/mess/config/configtest"
	"github.com/jostrzol/mess/pkg/board/boardtest"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/jostrzol/mess/pkg/mess/messtest"
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
				{color.White, "pawn", "D2"},
				{color.White, "pawn", "E2"},
				{color.White, "pawn", "F2"},
			},
			expectMoves: []move{
				{"E1", "C1"},
				{"E1", "G1"},
				{"E1", "D1"},
				{"E1", "F1"},
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

			matchers := make(map[string]messtest.MovesMatcherS, 0)
			for _, move := range tt.expectMoves {
				var matcher messtest.MovesMatcherS
				if oldMatcher, found := matchers[move.from]; found {
					matcher = oldMatcher
				} else {
					piece, err := game.Board().At(boardtest.NewSquare(move.from))
					assert.NoError(t, err)
					matcher.Piece = piece
				}
				matcher.Destinations = append(matcher.Destinations, move.to)
				matchers[move.from] = matcher
			}

			validMoves := game.ValidMoves()
			for from, matcher := range matchers {
				var foundMoves []mess.Move
				for _, validMove := range validMoves {
					if validMove.From.String() == from {
						foundMoves = append(foundMoves, validMove)
					}
				}

				messtest.MovesMatch(t, foundMoves, matcher)
			}
		})
	}
}
