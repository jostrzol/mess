package mess

import (
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/gen"
)

type Player struct {
	color     color.Color
	prisoners []*Piece
}

func NewPlayers() map[color.Color]*Player {
	colors := color.ColorValues()
	players := make(map[color.Color]*Player, len(colors))
	for _, color := range colors {
		players[color] = &Player{
			color: color,
		}
	}
	return players
}

func (p *Player) Color() color.Color {
	return p.color
}

func (p *Player) String() string {
	return p.color.String()
}

func (p *Player) Capture(piece *Piece) {
	p.prisoners = append(p.prisoners, piece)
}

func (p *Player) Prisoners() <-chan *Piece {
	return gen.Generator(p.prisoners)
}
