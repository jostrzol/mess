package player

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPlayers(t *testing.T) {
	players := NewPlayers()
	assert.Len(t, players, 2)

	for _, color := range ColorValues() {
		assert.Contains(t, players, color)
		player := players[color]
		assert.Equal(t, player.Color, color)
	}
}

func TestString(t *testing.T) {
	player := NewPlayers()[White]

	assert.Equal(t, player.String(), White.String())
}
