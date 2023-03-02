package player

import "github.com/jostrzol/mess/game/piece/color"

type Player struct {
	color color.Color
}

func (p *Player) Color() color.Color {
	return p.color
}

func (p *Player) String() string {
	return p.color.String()
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
