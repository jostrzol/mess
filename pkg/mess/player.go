package mess

import (
	"github.com/jostrzol/mess/pkg/board"
	brd "github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/event"
	"github.com/jostrzol/mess/pkg/utils"
)

type Player struct {
	color            color.Color
	pieces           map[*Piece]struct{}
	prisoners        map[*Piece]struct{}
	forwardDirection brd.Offset
}

func NewPlayers(board event.Subject) map[color.Color]*Player {
	players := map[color.Color]*Player{
		color.White: {
			forwardDirection: brd.Offset{X: 0, Y: 1},
		},
		color.Black: {
			forwardDirection: brd.Offset{X: 0, Y: -1},
		},
	}
	for color, player := range players {
		player.color = color
		player.pieces = make(map[*Piece]struct{})
		player.prisoners = make(map[*Piece]struct{})
		board.Observe(player)
	}
	return players
}

func (p *Player) Color() color.Color {
	return p.color
}

func (p *Player) Pieces() []*Piece {
	return utils.KeysToSlice(p.pieces)
}

func (p *Player) Prisoners() []*Piece {
	return utils.KeysToSlice(p.prisoners)
}

func (p *Player) ForwardDirection() board.Offset {
	return p.forwardDirection
}

func (p *Player) String() string {
	return p.color.String()
}

func (p *Player) Moves() []Move {
	result := make([]Move, 0)
	for piece := range p.pieces {
		result = append(result, piece.Moves()...)
	}
	return result
}

func (p *Player) AttackedSquares() []board.Square {
	resultSet := make(map[board.Square]struct{})
	for piece := range p.pieces {
		for _, move := range piece.Moves() {
			resultSet[move.To] = struct{}{}
		}
	}
	result := make([]board.Square, 0, len(resultSet))
	for square := range resultSet {
		result = append(result, square)
	}
	return result
}

func (p *Player) Handle(event event.Event) {
	switch e := event.(type) {
	case PiecePlaced:
		if e.Piece.Owner() == p {
			p.pieces[e.Piece] = struct{}{}
		}
		delete(p.prisoners, e.Piece)
	case PieceCaptured:
		if e.CapturedBy == p {
			p.prisoners[e.Piece] = struct{}{}
		}
	case PieceRemoved:
		if e.Piece.Owner() == p {
			delete(p.pieces, e.Piece)
		}
	}
}
