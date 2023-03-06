package mess

import (
	"testing"

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

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameSuite))
}
