package piece

import (
	"fmt"

	"github.com/jostrzol/mess/game/player"
)

type Type struct {
	Name string
}

func (pt *Type) String() string {
	return pt.Name
}

type Piece struct {
	Type  *Type
	Owner *player.Player
}

func (p *Piece) String() string {
	return fmt.Sprintf("%s %s", p.Owner.Color, p.Type)
}
