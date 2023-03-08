package mess

import (
	"github.com/jostrzol/mess/pkg/board"
	brd "github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/event"
	"github.com/jostrzol/mess/pkg/gen"
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

func (p *Player) Pieces() <-chan *Piece {
	return gen.FromKeys(p.pieces)
}

func (p *Player) Prisoners() <-chan *Piece {
	return gen.FromKeys(p.prisoners)
}

func (p *Player) ForwardDirection() board.Offset {
	return p.forwardDirection
}

func (p *Player) String() string {
	return p.color.String()
}

type Motion struct {
	Piece *Piece
	From  brd.Square
	To    brd.Square
}

func (m *Motion) Perform() {
	m.Piece.MoveTo(&m.To)
}

func (p *Player) GenerateMotions() []Motion {
	result := make([]Motion, 0)
	for piece := range p.Pieces() {
		for _, destination := range piece.GenerateMotions() {
			result = append(result, Motion{
				Piece: piece,
				From:  *piece.Square(),
				To:    destination,
			})
		}
	}
	return result
}

func (p *Player) Handle(event event.Event) {
	switch e := event.(type) {
	case PiecePlaced:
		if e.Piece.Owner() == p {
			p.pieces[e.Piece] = struct{}{}
		}
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
