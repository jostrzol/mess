package mess

import (
	"testing"

	"github.com/jostrzol/mess/pkg/board/boardtest"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/stretchr/testify/suite"
)

type GameSuite struct {
	suite.Suite
	game *State
}

func (s *GameSuite) SetupTest() {
	board, err := NewPieceBoard(8, 8)
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

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameSuite))
}
