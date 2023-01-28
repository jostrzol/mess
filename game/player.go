package game

//go:generate enumer -type=Color -transform=snake
type Color int

const (
	White Color = iota
	Black
)

type Player struct {
	Color Color
}

func (p *Player) String() string {
	return p.Color.String()
}

func NewPlayers() map[Color]*Player {
	colors := ColorValues()
	players := make(map[Color]*Player, len(colors))
	for _, color := range colors {
		players[color] = &Player{
			Color: color,
		}
	}
	return players
}
