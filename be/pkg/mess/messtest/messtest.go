package messtest

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/board/boardtest"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

type MovesMatcherS struct {
	Piece        *mess.Piece
	Destinations []string
}

func MovesMatcher(piece *mess.Piece, destinations ...string) MovesMatcherS {
	return MovesMatcherS{Piece: piece, Destinations: destinations}
}

func MovesMatch(t *testing.T, moves []*mess.MoveGroup, matchers ...MovesMatcherS) {
	anyNotFound := false
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("matching moves (%v) to (%v):\n", moves, matchers))

	for _, matcher := range matchers {
		for _, destination := range matcher.Destinations {
			dest := boardtest.NewSquare(destination)
			src := matcher.Piece.Square()
			found := -1
			for i, move := range moves {
				if move.Piece == matcher.Piece && move.From == src && move.To == dest {
					found = i
					break
				}
			}
			if found != -1 {
				moves[found] = moves[len(moves)-1]
				moves = moves[:len(moves)-1]
				msg.WriteString(fmt.Sprintf("FOUND:\t%v: %v->%v,\n", matcher.Piece, &src, &dest))
			} else {
				anyNotFound = true
				msg.WriteString(fmt.Sprintf("NOT FOUND:\t%v: %v->%v,\n", matcher.Piece, &src, &dest))
			}
		}
	}

	if len(moves) > 0 {
		for _, move := range moves {
			msg.WriteString(fmt.Sprintf("UNEXPECTED:\t%v,\n", move))
		}
	}

	if anyNotFound || len(moves) > 0 {
		t.Errorf(msg.String())
	}
}

func StaticMotion(t *testing.T, strings ...string) mess.Motion {
	t.Helper()
	return mess.Motion{
		Name: "test_generator",
		MoveGenerator: func(piece *mess.Piece) []board.Square {
			destinations := make([]board.Square, 0, len(strings))
			for _, squareStr := range strings {
				square, err := board.NewSquare(squareStr)
				assert.NoError(t, err)
				destinations = append(destinations, square)
			}
			return destinations
		},
	}
}

func OffsetMotion(t *testing.T, offsets ...board.Offset) mess.Motion {
	t.Helper()
	return mess.Motion{
		Name: "test_generator",
		MoveGenerator: func(piece *mess.Piece) []board.Square {
			destinations := make([]board.Square, 0, len(offsets))
			for _, offset := range offsets {
				square := piece.Square().Offset(offset)
				if piece.Board().Contains(square) {
					destinations = append(destinations, square)
				}
			}
			return destinations
		},
	}
}

func MoveWithOptionTexts(optionTexts []string, moveGroup *mess.MoveGroup) *mess.Move {
	moves := moveGroup.Moves()
	i := slices.IndexFunc(moves, func(move *mess.Move) bool {
		return fmt.Sprintf("%v", optionTexts) == fmt.Sprintf("%v", move.Options)
	})
	if i == -1 {
		err := fmt.Errorf("option set %v not generated", optionTexts)
		panic(err)
	}
	return moves[i]
}

func ChooseRandomRoute(src rand.Source, node *mess.OptionNode) mess.Route {
	visitor := &randomOptionDataVisitor{src: src}
	if node != nil {
		node.Accept(visitor)
	}
	return visitor.result
}

type randomOptionDataVisitor struct {
	src    rand.Source
	result mess.Route
}

func (r *randomOptionDataVisitor) VisitPieceTypeData(_ string, data mess.PieceTypeOptionData) {
	datum := getRandom(r.src, data.OptionData)
	r.handleDatum(datum)
}

func (r *randomOptionDataVisitor) VisitSquareData(_ string, data mess.SquareOptionData) {
	datum := getRandom(r.src, data.OptionData)
	r.handleDatum(datum)
}

func (r *randomOptionDataVisitor) VisitMoveData(_ string, data mess.MoveOptionData) {
	datum := getRandom(r.src, data.OptionData)
	r.handleDatum(datum)
}

func (r *randomOptionDataVisitor) VisitUnitData(_ string, data mess.UnitOptionData) {
	datum := getRandom(r.src, data.OptionData)
	r.handleDatum(datum)
}

func (r *randomOptionDataVisitor) handleDatum(datum mess.IOptionDatum) {
	r.result = append(r.result, datum.IOption())
	children := datum.NonEmptyChildren()
	if len(children) != 0 {
		nextNode := getRandom(r.src, children)
		nextNode.Accept(r)
	}
}

func getRandom[T any](src rand.Source, slice []T) T {
	index := int(src.Int63()) % len(slice)
	return slice[index]
}
