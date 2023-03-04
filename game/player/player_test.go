package player_test

import (
	"testing"

	"github.com/jostrzol/mess/game/piece/color"
	"github.com/jostrzol/mess/game/piece/piecetest"
	plr "github.com/jostrzol/mess/game/player"
	"github.com/jostrzol/mess/messtest/genassert"
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

func TestPrisonersEmpty(t *testing.T) {
	player := plr.NewPlayers()[color.White]

	genassert.Empty(t, player.Prisoners())
}

func TestPrisonersCapture(t *testing.T) {
	player := plr.NewPlayers()[color.White]
	knight := piecetest.Noones(piecetest.Knight(t))

	player.Capture(knight)

	genassert.Len(t, player.Prisoners(), 1)
	genassert.Contains(t, player.Prisoners(), knight)
}
