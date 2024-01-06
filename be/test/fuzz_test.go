package integration

import (
	"math/rand"
	"testing"

	"github.com/jostrzol/mess/pkg/mess/messtest"
	"github.com/jostrzol/mess/pkg/rules"
	"github.com/stretchr/testify/assert"
)

func FuzzChess(f *testing.F) {
	if testing.Short() {
		f.Skip()
	}

	fuzzGame(f, ChessRulesFile, 10)
}

func FuzzDobutsuShogi(f *testing.F) {
	fuzzGame(f, DobutsuShogiRulesFile, 50)
}

func FuzzHalma(f *testing.F) {
	fuzzGame(f, HalmaRulesFile, 50)
}

func fuzzGame(f *testing.F, rulesFilename string, maxSteps int) {
	f.Helper()

	f.Add(int64(12345))
	f.Fuzz(func(t *testing.T, seed int64) {
		game, err := rules.DecodeRulesFromOs(rulesFilename, true)
		assert.NoError(t, err)

		src := rand.NewSource(seed)
		src.Int63()

		resolution := game.Resolution()
		for i := 0; !resolution.DidEnd && i < maxSteps; i++ {
			options, err := game.TurnOptions()
			assert.NoError(t, err)

			route := messtest.ChooseRandomRoute(src, options)
			err = game.PlayTurn(route)
			assert.NoError(t, err)

			resolution = game.Resolution()
		}
	})
}
