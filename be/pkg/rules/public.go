package rules

import (
	"fmt"
	"os"

	"github.com/jostrzol/mess/pkg/mess"
)

type File struct {
	Src      []byte
	Filename string
}

func DecodeRulesFromOs(filename string, placePieces bool) (*mess.Game, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("opening rules file: %w", err)
	}

	return DecodeRules(&File{src, filename}, placePieces)
}

func DecodeRules(file *File, placePieces bool) (*mess.Game, error) {
	ctx := InitialEvalContext

	rules, err := decodeRules(file.Src, file.Filename, ctx)
	if err != nil {
		return nil, fmt.Errorf("decoding rules: %w", err)
	}

	game, err := rules.toEmptyGameState(ctx)
	if err != nil {
		return nil, fmt.Errorf("initializing game from rules: %w", err)
	}

	if placePieces {
		err = rules.placePieces(game.State)
		if err != nil {
			return nil, fmt.Errorf("placing initial pieces: %w", err)
		}
	}

	return game, nil
}
