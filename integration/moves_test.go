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
	// blank to end turn
	type move struct {
		from string
		to   string
	}
	type test struct {
		name        string
		initState   []piece
		whenMoves   []move
		expectMoves map[string][]string
	}

	tests := []test{
		{
			name: "castling",
			initState: []piece{
				{color.White, "king", "E1"},
				{color.White, "rook", "A1"},
				{color.White, "rook", "H1"},
			},
			expectMoves: map[string][]string{
				"E1": {"C1", "G1", "D1", "F1", "D2", "E2", "F2"},
			},
		},
		{
			name: "castling_enemy_3_and_further_from_king",
			initState: []piece{
				{color.White, "king", "E1"},
				{color.White, "rook", "A1"},
				{color.White, "rook", "H1"},
				{color.Black, "rook", "A3"},
				{color.Black, "rook", "B3"},
				{color.Black, "rook", "H3"},
			},
			expectMoves: map[string][]string{
				"E1": {"C1", "G1", "D1", "F1", "D2", "E2", "F2"},
			},
		},
		{
			name: "castling_blocked_by_own_knights",
			initState: []piece{
				{color.White, "king", "E1"},
				{color.White, "rook", "A1"},
				{color.White, "rook", "H1"},
				{color.White, "knight", "B1"},
				{color.White, "knight", "G1"},
			},
			expectMoves: map[string][]string{
				"E1": {"D1", "F1", "D2", "E2", "F2"},
			},
		},
		{
			name: "castling_blocked_by_own_bishops",
			initState: []piece{
				{color.White, "king", "E1"},
				{color.White, "rook", "A1"},
				{color.White, "rook", "H1"},
				{color.White, "bishop", "C1"},
				{color.White, "bishop", "F1"},
			},
			expectMoves: map[string][]string{
				"E1": {"D1", "D2", "E2", "F2"},
			},
		},
		{
			name: "castling_long_blocked_by_own_queen",
			initState: []piece{
				{color.White, "king", "E1"},
				{color.White, "rook", "A1"},
				{color.White, "rook", "H1"},
				{color.White, "queen", "D1"},
			},
			expectMoves: map[string][]string{
				"E1": {"G1", "F1", "D2", "E2", "F2"},
			},
		},
		{
			name: "castling_blocked_by_enemy_on_king",
			initState: []piece{
				{color.White, "king", "E1"},
				{color.White, "rook", "A1"},
				{color.White, "rook", "H1"},
				{color.Black, "rook", "E3"},
			},
			expectMoves: map[string][]string{
				"E1": {"D1", "F1", "D2", "F2"},
			},
		},
		{
			name: "castling_blocked_by_enemy_1_from_king",
			initState: []piece{
				{color.White, "king", "E1"},
				{color.White, "rook", "A1"},
				{color.White, "rook", "H1"},
				{color.Black, "rook", "D3"},
				{color.Black, "rook", "F3"},
			},
			expectMoves: map[string][]string{
				"E1": {"E2"},
			},
		},
		{
			name: "castling_blocked_by_enemy_2_from_king",
			initState: []piece{
				{color.White, "king", "E1"},
				{color.White, "rook", "A1"},
				{color.White, "rook", "H1"},
				{color.Black, "rook", "C3"},
				{color.Black, "rook", "G3"},
			},
			expectMoves: map[string][]string{
				"E1": {"D1", "F1", "D2", "E2", "F2"},
			},
		},
		{
			name: "castling_blocked_after_king_move",
			initState: []piece{
				{color.White, "king", "E1"},
				{color.White, "rook", "A1"},
				{color.White, "rook", "H1"},
			},
			whenMoves: []move{{"E1", "E2"}, {"E2", "E1"}},
			expectMoves: map[string][]string{
				"E1": {"D1", "F1", "D2", "E2", "F2"},
			},
		},
		{
			name: "castling_long_blocked_after_rook_move",
			initState: []piece{
				{color.White, "king", "E1"},
				{color.White, "rook", "A1"},
				{color.White, "rook", "H1"},
			},
			whenMoves: []move{{"A1", "A2"}, {"A2", "A1"}},
			expectMoves: map[string][]string{
				"E1": {"G1", "D1", "F1", "D2", "E2", "F2"},
			},
		},
		{
			name: "castling_short_blocked_after_rook_move",
			initState: []piece{
				{color.White, "king", "E1"},
				{color.White, "rook", "A1"},
				{color.White, "rook", "H1"},
			},
			whenMoves: []move{{"H1", "H2"}, {"H2", "H1"}},
			expectMoves: map[string][]string{
				"E1": {"C1", "D1", "F1", "D2", "E2", "F2"},
			},
		},
		{
			name: "en_passant_left",
			initState: []piece{
				{color.White, "pawn", "A4"},
				{color.Black, "pawn", "B7"},
			},
			whenMoves: []move{{"A4", "A5"}, {"B7", "B5"}},
			expectMoves: map[string][]string{
				"A5": {"A6", "B6"},
			},
		},
		{
			name: "en_passant_right",
			initState: []piece{
				{color.White, "pawn", "C4"},
				{color.Black, "pawn", "B7"},
			},
			whenMoves: []move{{"C4", "C5"}, {"B7", "B5"}},
			expectMoves: map[string][]string{
				"C5": {"C6", "B6"},
			},
		},
		{
			name: "en_passant_blocked_if_not_performed_asap",
			initState: []piece{
				{color.White, "pawn", "A4"},
				{color.Black, "pawn", "B7"},
				{color.White, "pawn", "G2"},
				{color.Black, "pawn", "G7"},
			},
			whenMoves: []move{{"A4", "A5"}, {"B7", "B5"}, {"G2", "G3"}, {"G7", "G6"}},
			expectMoves: map[string][]string{
				"A5": {"A6"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game, err := config.DecodeConfig(RulesFile, configtest.PanicInteractor{}, false)
			assert.NoError(t, err)

			for _, piece := range tt.initState {
				place(t, game, piece.color, piece.name, piece.square)
			}

			for _, move := range tt.whenMoves {
				piece, err := game.Board().At(boardtest.NewSquare(move.from))
				assert.NoError(t, err)
				if piece.Owner() != game.CurrentPlayer() {
					game.EndTurn()
				}

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
			for from, to := range tt.expectMoves {
				var matcher messtest.MovesMatcherS
				if oldMatcher, found := matchers[from]; found {
					matcher = oldMatcher
				} else {
					piece, err := game.Board().At(boardtest.NewSquare(from))
					assert.NoError(t, err)
					matcher.Piece = piece
				}
				matcher.Destinations = append(matcher.Destinations, to...)
				matchers[from] = matcher

				piece, _ := game.Board().At(boardtest.NewSquare(from))
				if piece != nil && piece.Owner() != game.CurrentPlayer() {
					game.EndTurn()
				}
				println(game.CurrentPlayer().Color().String())
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
