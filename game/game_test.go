package game_test

import (
	"testing"

	"github.com/jostrzol/mess/game"
	"github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/board/boardtest"
	"github.com/jostrzol/mess/game/piece"
	"github.com/jostrzol/mess/game/piece/piecetest"
	"github.com/jostrzol/mess/game/player"
	"github.com/stretchr/testify/suite"
)

type GameSuite struct {
	suite.Suite
	game *game.State
}

func (s *GameSuite) SetupTest() {
	players := player.NewPlayers()
	board, err := board.NewBoard[*piece.Piece](8, 8)
	s.NoError(err)

	s.game = &game.State{
		Board:   board,
		Players: players,
	}
}

func (s *GameSuite) TestGetPlayer() {
	for _, color := range player.ColorValues() {
		s.Run(color.String(), func() {
			player, err := s.game.GetPlayer(color)
			s.NoError(err)
			s.Equal(player.Color, color)
		})
	}
}

func (s *GameSuite) TestGetPlayerNotFound() {
	_, err := s.game.GetPlayer(player.Color(-1))
	s.Error(err)
}

func (s *GameSuite) TestPiecesPerPlayer() {
	t := s.T()

	white := s.game.Players[player.White]
	black := s.game.Players[player.Black]

	rookW := &piece.Piece{Type: piecetest.Rook(t), Owner: white}
	knightW := &piece.Piece{Type: piecetest.Knight(t), Owner: white}
	rookB := &piece.Piece{Type: piecetest.Rook(t), Owner: black}

	s.game.Board.Place(rookW, boardtest.NewSquare("A1"))
	s.game.Board.Place(knightW, boardtest.NewSquare("B4"))
	s.game.Board.Place(rookB, boardtest.NewSquare("F2"))

	results := s.game.PiecesPerPlayer()
	s.Len(results, 2)

	s.ElementsMatch(results[white], []*piece.Piece{rookW, knightW})
	s.ElementsMatch(results[black], []*piece.Piece{rookB})
}

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameSuite))
}
