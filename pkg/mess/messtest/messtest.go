package messtest

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/board/boardtest"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/stretchr/testify/assert"
)

type MovesMatcherS struct {
	Piece        *mess.Piece
	Destinations []string
}

func MovesMatcher(piece *mess.Piece, destinations ...string) MovesMatcherS {
	return MovesMatcherS{Piece: piece, Destinations: destinations}
}

func MovesMatch(t *testing.T, moves []mess.Move, matchers ...MovesMatcherS) {
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
			msg.WriteString(fmt.Sprintf("UNEXPECTED:\t%v,\n", &move))
		}
	}

	if anyNotFound || len(moves) > 0 {
		t.Errorf(msg.String())
	}
}

func StaticMoveGenerator(t *testing.T, strings ...string) mess.MoveGenerator {
	t.Helper()
	return func(piece *mess.Piece) ([]board.Square, mess.MoveAction) {
		destinations := make([]board.Square, 0, len(strings))
		for _, squareStr := range strings {
			square, err := board.NewSquare(squareStr)
			assert.NoError(t, err)
			destinations = append(destinations, square)
		}
		return destinations, nil
	}
}

func OffsetMoveGenerator(t *testing.T, offsets ...board.Offset) mess.MoveGenerator {
	t.Helper()
	return func(piece *mess.Piece) ([]board.Square, mess.MoveAction) {
		destinations := make([]board.Square, 0, len(offsets))
		for _, offset := range offsets {
			square := piece.Square().Offset(offset)
			if piece.Board().Contains(square) {
				destinations = append(destinations, square)
			}
		}
		return destinations, nil
	}
}
