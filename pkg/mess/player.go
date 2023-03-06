package mess

import (
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/event"
	"github.com/jostrzol/mess/pkg/gen"
)

type Player struct {
	color     color.Color
	prisoners []*Piece
}

func NewPlayers(board event.Subject) map[color.Color]*Player {
	colors := color.ColorValues()
	players := make(map[color.Color]*Player, len(colors))
	for _, color := range colors {
		player := &Player{color: color}
		board.Observe(player)
		players[color] = player
	}
	return players
}

func (p *Player) Color() color.Color {
	return p.color
}

func (p *Player) String() string {
	return p.color.String()
}

func (p *Player) Handle(event event.Event) {
	switch e := event.(type) {
	case PieceCaptured:
		if e.CapturedBy == p {
			p.prisoners = append(p.prisoners, e.Piece)
		}
	}
}

func (p *Player) Prisoners() <-chan *Piece {
	return gen.Generator(p.prisoners)
}
