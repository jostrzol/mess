package rules

import (
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/hashicorp/hcl/v2"
	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/zclconf/go-cty/cty"
)

func (c *rules) toEmptyGameState(ctx *hcl.EvalContext) (*mess.Game, error) {
	brd, err := mess.NewPieceBoard(int(c.Board.Width), int(c.Board.Height))
	if err != nil {
		return nil, fmt.Errorf("creating new board: %w", err)
	}

	state := mess.NewState(brd)
	controller := newController(state, ctx, c)

	game := mess.NewGame(state, controller)

	stateValidators, err := controller.GetStateValidators()
	if err != nil {
		return nil, fmt.Errorf("parsing state validators: %w", err)
	}
	for _, validator := range stateValidators {
		state.AddStateValidator(validator)
	}

	for _, pieceTypeRules := range c.PieceTypes.PieceTypes {
		pieceType, err := decodePieceType(controller, pieceTypeRules)
		if err != nil {
			return nil, fmt.Errorf("decoding piece type %q: %v", pieceTypeRules.Name, err)
		}
		state.AddPieceType(pieceType)
	}

	assets, err := decodeAssets(c.Assets)
	if err != nil {
		return nil, err
	}
	state.Assets = assets

	initializeContext(ctx, game)
	return game, nil
}

func decodePieceType(controller *controller, pieceTypeRules pieceTypeRules) (*mess.PieceType, error) {
	pieceType := mess.NewPieceType(pieceTypeRules.Name)
	for _, motionRules := range pieceTypeRules.Motions {
		moveGenerator, err := controller.GetCustomFuncAsGenerator(motionRules.GeneratorName)
		if err != nil {
			return nil, err
		}
		var action mess.MoveActionFunc
		if motionRules.ActionName != "" {
			action, err = controller.GetCustomFuncAsAction(motionRules.ActionName)
			if err != nil {
				return nil, err
			}
		}
		var choiceFunction mess.MoveChoiceFunc
		if motionRules.ChoiceFunctionName != "" {
			choiceFunction, err = controller.GetCustomFuncAsChoiceFunction(
				motionRules.ChoiceFunctionName)
			if err != nil {
				return nil, err
			}
		}
		pieceType.AddMotion(
			mess.Motion{
				Name:          motionRules.GeneratorName,
				MoveGenerator: moveGenerator,
				ChoiceFunc:    choiceFunction,
				Action:        action,
			},
		)
	}
	if pieceTypeRules.Representation != nil {
		if pieceTypeRules.Representation.Black != nil {
			representation, err := decodeRepresentation(pieceTypeRules.Representation.Black)
			if err != nil {
				return nil, fmt.Errorf("decoding representation: %w", err)
			}
			pieceType.SetRepresentation(color.Black, representation)
		}
		if pieceTypeRules.Representation.White != nil {
			representation, err := decodeRepresentation(pieceTypeRules.Representation.White)
			if err != nil {
				return nil, fmt.Errorf("decoding representation: %w", err)
			}
			pieceType.SetRepresentation(color.White, representation)
		}
	}
	return pieceType, nil
}

func decodeRepresentation(representation *representation) (mess.Representation, error) {
	var symbol rune
	var icon mess.AssetKey
	var err error
	if representation.Symbol != nil {
		symbol, err = decodeSymbol(*representation.Symbol)
		if err != nil {
			return mess.Representation{}, fmt.Errorf("decoding symbol: %w", err)
		}
	}
	if representation.Icon != nil {
		icon = mess.NewAssetKey(*representation.Icon)
	}
	return mess.Representation{Symbol: symbol, Icon: icon}, err
}

func decodeSymbol(symbol string) (rune, error) {
	r, n := utf8.DecodeRuneInString(symbol)
	if n == 0 {
		return 0, fmt.Errorf("symbol cannot be empty")
	} else if r == utf8.RuneError {
		return 0, fmt.Errorf("symbol not an utf-8 character")
	} else if n != len(symbol) {
		return 0, fmt.Errorf("symbol too long (must be exactly one utf-8 character)")
	}
	return r, nil
}

func (c *rules) placePieces(state *mess.State) error {
	placementRules := map[color.Color]map[string]string{
		color.White: c.InitialState.WhitePieces,
		color.Black: c.InitialState.BlackPieces,
	}
	for color, pieces := range placementRules {
		player := state.Player(color)

		for squareString, pieceTypeName := range pieces {
			square, err := board.NewSquare(squareString)
			if err != nil {
				return fmt.Errorf("parsing square: %w", err)
			}

			pieceType, err := state.GetPieceType(pieceTypeName)
			if err != nil {
				return fmt.Errorf("getting piece type: %w", err)
			}

			piece := mess.NewPiece(pieceType, player)

			err = piece.PlaceOn(state.Board(), square)
			if err != nil {
				return fmt.Errorf("placing a piece: %w", err)
			}
		}
	}

	return nil
}

func decodeAssets(assetsCty cty.Value) (mess.Assets, error) {
	result := make(mess.Assets)
	type keyAssetPair struct {
		key   string
		value cty.Value
	}
	assets := []keyAssetPair{{"", assetsCty}}
	for len(assets) != 0 {
		asset := assets[0]
		assets = assets[1:]

		if asset.value.Type().IsObjectType() {
			for key, value := range asset.value.AsValueMap() {
				assets = append(assets, keyAssetPair{asset.key + "/" + key, value})
			}
		} else {
			valueRaw := asset.value.AsString()
			whitespaceReplacer := strings.NewReplacer(" ", "", "\t", "", "\n", "", "\r", "")
			valueB64 := whitespaceReplacer.Replace(valueRaw)
			b64Reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(valueB64))
			gzipReader, err := gzip.NewReader(b64Reader)
			if err != nil {
				return nil, fmt.Errorf("decoding asset %v: %w", asset.key, err)
			}
			value, err := io.ReadAll(gzipReader)
			if err != nil {
				return nil, fmt.Errorf("decoding asset %v: %w", asset.key, err)
			}
			result[mess.NewAssetKey(asset.key)] = value
		}
	}
	return result, nil
}
