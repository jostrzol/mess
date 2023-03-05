package mess

import (
	"testing"

	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/board/boardtest"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/genassert"
	"github.com/stretchr/testify/suite"
)

type GameSuite struct {
	suite.Suite
	game *State
}

func (s *GameSuite) SetupTest() {
	board, err := board.NewBoard[*Piece](8, 8)
	s.NoError(err)
	s.game = NewState(board)
}

func (s *GameSuite) TestGetPlayer() {
	for _, color := range color.ColorValues() {
		s.Run(color.String(), func() {
			player := s.game.Player(color)
			s.Equal(player.Color(), color)
		})
	}
}

func (s *GameSuite) TestGetPlayerNotFound() {
	s.Panics(func() {
		s.game.Player(color.Color(-1))
	})
}

func (s *GameSuite) TestPiecesPerPlayer() {
	t := s.T()

	white := s.game.Players[color.White]
	black := s.game.Players[color.Black]

	rookW := NewPiece(Rook(t), white)
	knightW := NewPiece(Knight(t), white)
	rookB := NewPiece(Rook(t), black)

	rookW.PlaceOn(s.game.Board, boardtest.NewSquare("A1"))
	knightW.PlaceOn(s.game.Board, boardtest.NewSquare("B4"))
	rookB.PlaceOn(s.game.Board, boardtest.NewSquare("F2"))

	results := s.game.PiecesPerPlayer()
	s.Len(results, 2)

	s.ElementsMatch(results[white], []*Piece{rookW, knightW})
	s.ElementsMatch(results[black], []*Piece{rookB})
}

func (s *GameSuite) TestMoveNoCapture() {
	white := s.game.Players[color.White]
	rook := NewPiece(Rook(s.T()), white)
	rook.PlaceOn(s.game.Board, boardtest.NewSquare("A1"))

	err := s.game.Move(rook, boardtest.NewSquare("A2"))
	s.NoError(err)

	genassert.Empty(s.T(), white.Prisoners())
}

func (s *GameSuite) TestMoveCapture() {
	white := s.game.Players[color.White]
	black := s.game.Players[color.Black]

	knight := NewPiece(Knight(s.T()), white)
	knight.PlaceOn(s.game.Board, boardtest.NewSquare("A2"))
	rook := NewPiece(Rook(s.T()), black)
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
