package integration

import (
	"math/rand"
	"testing"

	"github.com/jostrzol/mess/pkg/mess/messtest"
	"github.com/jostrzol/mess/pkg/rules"
	"github.com/stretchr/testify/assert"
)

func FuzzGameMax10Steps(f *testing.F) {
	if testing.Short() {
		f.Skip()
	}

	f.Add(int64(12345))
	f.Fuzz(func(t *testing.T, seed int64) {
		game, err := rules.DecodeRulesFromOs(ChessRulesFile, true)
		assert.NoError(t, err)

		src := rand.NewSource(seed)
		src.Int63()

		resolution := game.Resolution()
		for i := 0; !resolution.DidEnd && i < 10; i++ {
			moves := game.ValidMoves()
			assert.NotEmpty(t, moves)

			chosen := int(src.Int63()) % len(moves)
			moveGroup := moves[chosen]

			move := messtest.ChooseRandom(src, moveGroup)
			err = move.Perform()
			assert.NoError(t, err)

			resolution = game.Resolution()
		}
	})
}
