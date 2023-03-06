package mess

import (
	"testing"

	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/event"
	"github.com/jostrzol/mess/pkg/genassert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestNewPlayers(t *testing.T) {
	board := event.NewSubject()

	players := NewPlayers(board)
	assert.Len(t, players, 2)

	for _, color := range color.ColorValues() {
		assert.Contains(t, players, color)
		player := players[color]
		assert.Equal(t, player.Color(), color)
		assert.True(t, board.IsObserving(player))
	}
}

type PlayerSuite struct {
	suite.Suite
	board   event.Subject
	players map[color.Color]*Player
	white   *Player
	black   *Player
}

func (s *PlayerSuite) SetupTest() {
	s.board = event.NewSubject()
	s.players = NewPlayers(s.board)
	s.white = s.players[color.White]
	s.black = s.players[color.Black]
}

func (s *PlayerSuite) TestString() {
	s.Equal(s.white.String(), color.White.String())
}

func (s *PlayerSuite) TestPrisonersEmpty() {
	genassert.Empty(s.T(), s.white.Prisoners())
}

func (s *PlayerSuite) TestPrisonersCapture() {
	knight := NewPiece(Knight(s.T()), s.black)
	s.board.Notify(PieceCaptured{
		Piece:      knight,
		CapturedBy: s.white,
	})

	genassert.Len(s.T(), s.white.Prisoners(), 1)
	genassert.Contains(s.T(), s.white.Prisoners(), knight)

	genassert.Empty(s.T(), s.black.Prisoners())
}

func TestPlayerSource(t *testing.T) {
	suite.Run(t, new(PlayerSuite))
}
