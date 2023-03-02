package player_test

import (
	"testing"

	"github.com/jostrzol/mess/game/piece/color"
	plr "github.com/jostrzol/mess/game/player"
	"github.com/stretchr/testify/assert"
)

func TestNewPlayers(t *testing.T) {
	players := plr.NewPlayers()
	assert.Len(t, players, 2)

	for _, color := range color.ColorValues() {
		assert.Contains(t, players, color)
		player := players[color]
		assert.Equal(t, player.Color(), color)
	}
}

func TestString(t *testing.T) {
	player := plr.NewPlayers()[color.White]

	assert.Equal(t, player.String(), color.White.String())
}
