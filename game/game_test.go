package game_test

import (
	"testing"

	"github.com/jostrzol/mess/game"
	"github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/board/boardtest"
	"github.com/jostrzol/mess/game/piece"
	"github.com/jostrzol/mess/game/piece/color"
	"github.com/jostrzol/mess/game/piece/piecetest"
	"github.com/jostrzol/mess/messtest/genassert"
	"github.com/stretchr/testify/suite"
)

type GameSuite struct {
	suite.Suite
	game *game.State
}

func (s *GameSuite) SetupTest() {
	board, err := board.NewBoard[*piece.Piece](8, 8)
	s.NoError(err)
	s.game = game.NewState(board)
}

func (s *GameSuite) TestGetPlayer() {
	for _, color := range color.ColorValues() {
		s.Run(color.String(), func() {
			player := s.game.GetPlayer(color)
			s.Equal(player.Color(), color)
		})
	}
}

func (s *GameSuite) TestGetPlayerNotFound() {
	s.Panics(func() {
		s.game.GetPlayer(color.Color(-1))
	})
}

func (s *GameSuite) TestPiecesPerPlayer() {
	t := s.T()

	white := s.game.Players[color.White]
	black := s.game.Players[color.Black]

	rookW := piece.NewPiece(piecetest.Rook(t), white)
	knightW := piece.NewPiece(piecetest.Knight(t), white)
	rookB := piece.NewPiece(piecetest.Rook(t), black)

	rookW.PlaceOn(s.game.Board, boardtest.NewSquare("A1"))
	knightW.PlaceOn(s.game.Board, boardtest.NewSquare("B4"))
	rookB.PlaceOn(s.game.Board, boardtest.NewSquare("F2"))

	results := s.game.PiecesPerPlayer()
	s.Len(results, 2)

	s.ElementsMatch(results[white], []*piece.Piece{rookW, knightW})
	s.ElementsMatch(results[black], []*piece.Piece{rookB})
}

func (s *GameSuite) TestMoveNoCapture() {
	white := s.game.Players[color.White]
	rook := piece.NewPiece(piecetest.Rook(s.T()), white)
	rook.PlaceOn(s.game.Board, boardtest.NewSquare("A1"))

	err := s.game.Move(rook, boardtest.NewSquare("A2"))
	s.NoError(err)

	genassert.Empty(s.T(), white.Prisoners())
}

func (s *GameSuite) TestMoveCapture() {
	white := s.game.Players[color.White]
	black := s.game.Players[color.Black]

	knight := piece.NewPiece(piecetest.Knight(s.T()), white)
	knight.PlaceOn(s.game.Board, boardtest.NewSquare("A2"))
	rook := piece.NewPiece(piecetest.Rook(s.T()), black)
	rook.PlaceOn(s.game.Board, boardtest.NewSquare("A1"))

	err := s.game.Move(rook, boardtest.NewSquare("A2"))
	s.NoError(err)

	genassert.Empty(s.T(), white.Prisoners())
	genassert.Len(s.T(), black.Prisoners(), 1)
	genassert.Contains(s.T(), black.Prisoners(), knight)
}

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameSuite))
}
