package integration

import (
	"math/rand"
	"testing"

	"github.com/jostrzol/mess/config"
	"github.com/jostrzol/mess/config/configtest"
	"github.com/stretchr/testify/assert"
)

func FuzzGameMax10Steps(f *testing.F) {
	if testing.Short() {
		f.Skip()
	}

	f.Add(int64(12345))
	f.Fuzz(func(t *testing.T, seed int64) {
		game, err := config.DecodeConfig("../rules.hcl", configtest.RandomInteractor{})
		assert.NoError(t, err)

		src := rand.NewSource(seed)
		src.Int63()

		isFinished, _ := game.PickWinner()
		for i := 0; !isFinished && i < 10; i++ {
			moves := game.ValidMoves()
			assert.NotEmpty(t, moves)

			chosen := int(src.Int63()) % len(moves)
			err := moves[chosen].Perform()
			assert.NoError(t, err)

			isFinished, _ = game.PickWinner()
		}
	})
}
