package mess

import (
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/event"
	"github.com/jostrzol/mess/pkg/gen"
)

type Player struct {
	color     color.Color
	pieces    map[*Piece]struct{}
	prisoners map[*Piece]struct{}
}

func NewPlayers(board event.Subject) map[color.Color]*Player {
	colors := color.ColorValues()
	players := make(map[color.Color]*Player, len(colors))
	for _, color := range colors {
		player := &Player{
			color:     color,
			pieces:    make(map[*Piece]struct{}),
			prisoners: make(map[*Piece]struct{}),
		}
		board.Observe(player)
		players[color] = player
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

func (p *Player) String() string {
	return p.color.String()
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
