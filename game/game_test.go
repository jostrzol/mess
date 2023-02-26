package game

import (
	"testing"

	"github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/piece"
	"github.com/jostrzol/mess/game/player"
	"github.com/stretchr/testify/suite"
)

type GameSuite struct {
	suite.Suite
	game *State
}

func (s *GameSuite) SetupTest() {
	players := player.NewPlayers()
	board, err := board.NewBoard(8, 8)
	s.NoError(err)

	s.game = &State{
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
	square1, _ := board.NewSquare("A1")
	square2, _ := board.NewSquare("B4")
	square3, _ := board.NewSquare("F2")

	white := s.game.Players[player.White]
	black := s.game.Players[player.Black]

	rookW := &piece.Piece{Type: piece.Rook(), Owner: white}
	knightW := &piece.Piece{Type: piece.Knight(), Owner: white}
	rookB := &piece.Piece{Type: piece.Rook(), Owner: black}

	s.game.Board.Place(rookW, square1)
	s.game.Board.Place(knightW, square2)
	s.game.Board.Place(rookB, square3)

	results := s.game.PiecesPerPlayer()
	s.Len(results, 2)

	s.ElementsMatch(results[white], []board.PieceOnSquare{
		{Piece: rookW, Square: square1},
		{Piece: knightW, Square: square2},
	})
	s.ElementsMatch(results[black], []board.PieceOnSquare{
		{Piece: rookB, Square: square3},
	})
}

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameSuite))
}
