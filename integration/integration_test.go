package integration

import (
	"math/rand"
	"testing"

	"github.com/jostrzol/mess/config"
	"github.com/stretchr/testify/assert"
)

func FuzzGameMax10Steps(f *testing.F) {
	if testing.Short() {
		f.Skip()
	}

	f.Add(int64(12345))
	f.Fuzz(func(t *testing.T, seed int64) {
		state, controller, err := config.DecodeConfig("../rules.hcl")
		assert.NoError(t, err)

		src := rand.NewSource(seed)
		src.Int63()

		isFinished, _ := controller.PickWinner(state)
		for i := 0; !isFinished && i < 10; i++ {
			moves := state.ValidMoves()
			assert.NotEmpty(t, moves)

			chosen := int(src.Int63()) % len(moves)
			err := moves[chosen].Perform()
			assert.NoError(t, err)

			isFinished, _ = controller.PickWinner(state)
		}
	})
}
