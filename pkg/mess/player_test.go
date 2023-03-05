package mess

import (
	"testing"

	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/genassert"
	"github.com/stretchr/testify/assert"
)

func TestNewPlayers(t *testing.T) {
	players := NewPlayers()
	assert.Len(t, players, 2)

	for _, color := range color.ColorValues() {
		assert.Contains(t, players, color)
		player := players[color]
		assert.Equal(t, player.Color(), color)
	}
}

func TestString(t *testing.T) {
	player := NewPlayers()[color.White]

	assert.Equal(t, player.String(), color.White.String())
}

func TestPrisonersEmpty(t *testing.T) {
	player := NewPlayers()[color.White]

	genassert.Empty(t, player.Prisoners())
}

func TestPrisonersCapture(t *testing.T) {
	player := NewPlayers()[color.White]
	knight := Noones(Knight(t))

	player.Capture(knight)

	genassert.Len(t, player.Prisoners(), 1)
	genassert.Contains(t, player.Prisoners(), knight)
}
