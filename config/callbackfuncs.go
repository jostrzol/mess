package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/jostrzol/mess/game"
	"github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/piece"
	"github.com/jostrzol/mess/game/piece/color"
	plr "github.com/jostrzol/mess/game/player"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

type callbackFunctionsConfig struct {
	DecideWinnerFunc function.Function            `mapstructure:"decide_winner"`
	CustomFuncs      map[string]function.Function `mapstructure:",remain"`
}

func (c *callbackFunctionsConfig) DecideWinner(state *game.State) *plr.Player {
	ctyState := gameStateToCty(state)
	ctyWinner, err := c.DecideWinnerFunc.Call([]cty.Value{ctyState})
	if err != nil {
		log.Printf("calling user-defined function: %v", err)
		return nil
	}
	winner, err := playerFromCty(state, ctyWinner)
	if err != nil {
		log.Printf("getting winner: %v", err)
		return nil
	}
	return winner
}

func (c *callbackFunctionsConfig) GetCustomFuncAsGenerator(name string) (piece.MotionGenerator, error) {
	funcCty, ok := c.CustomFuncs[name]
	if !ok {
		return nil, fmt.Errorf("user function %q not found", name)
	}

	return piece.FuncMotionGenerator(func(piece *piece.Piece) []board.Square {
		pieceCty := pieceToCty(piece)
		squareCty := squareToCty(&piece.Square)
		result, err := funcCty.Call([]cty.Value{squareCty, pieceCty})
		if err != nil {
			log.Printf("calling motion generator for %v at %v: %v", piece, piece.Square, err)
			return make([]board.Square, 0)
		}

		squares, err := squaresFromCty(&result)
		if err != nil {
			log.Printf("parsing motion generator result: %v", err)
		}
		return squares
	}), nil
}

func gameStateToCty(state *game.State) cty.Value {
	piecesPerPlayer := state.PiecesPerPlayer()
	players := make(map[string]cty.Value, len(state.Players))
	for _, player := range state.Players {
		pieces := piecesPerPlayer[player]
		players[player.Color().String()] = playerToCty(player, pieces)
	}
	return cty.ObjectVal(map[string]cty.Value{
		"players": cty.MapVal(players),
	})
}

func playerToCty(player *plr.Player, pieces []*piece.Piece) cty.Value {
	piecesCty := make([]cty.Value, len(pieces))
	for i, piece := range pieces {
		piecesCty[i] = pieceToCty(piece)
	}
	return cty.ObjectVal(map[string]cty.Value{
		"color":  cty.StringVal(player.Color().String()),
		"pieces": cty.ListVal(piecesCty),
	})
}

func pieceToCty(piece *piece.Piece) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"type":   cty.StringVal(piece.Type.Name),
		"square": cty.StringVal(piece.Square.String()),
	})
}

func squareToCty(square *board.Square) cty.Value {
	return cty.StringVal(square.String())
}

type manyErrors []error

func (errors manyErrors) Error() string {
	var b strings.Builder
	b.WriteString("[\n")
	for _, err := range errors {
		b.WriteByte('\t')
		b.WriteString(err.Error())
		b.WriteByte('\n')
	}
	return b.String()
}

func squaresFromCty(value *cty.Value) ([]board.Square, error) {
	ty := value.Type()
	if ty != cty.List(cty.String) {
		err := fmt.Errorf("expected type %v, got %v", cty.List(cty.String), ty)
		return make([]board.Square, 0), err
	}
	errors := make(manyErrors, 0)
	destinations := make([]board.Square, value.LengthInt())
	for _, squareCty := range value.AsValueSlice() {
		squareStr := squareCty.AsString()
		square, err := board.NewSquare(squareStr)
		if err != nil {
			errors = append(errors, fmt.Errorf("parsing square %q: %w", squareStr, err))
		} else {
			destinations = append(destinations, *square)
		}
	}
	if len(errors) > 0 {
		return destinations, errors
	}
	return destinations, nil
}

func playerFromCty(state *game.State, player cty.Value) (*plr.Player, error) {
	if player.IsNull() {
		return nil, nil
	}
	winnerColorStr := player.GetAttr("color").AsString()
	winnerColor, err := color.ColorString(winnerColorStr)
	if err != nil {
		return nil, fmt.Errorf("parsing player color: %w", err)
	}
	winner := state.GetPlayer(winnerColor)
	return winner, nil
}
