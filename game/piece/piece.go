package piece

import (
	"fmt"

	"github.com/jostrzol/mess/game/player"
)

type PieceType struct {
	Name string
}

func (pt *PieceType) String() string {
	return pt.Name
}

type Piece struct {
	Type  *PieceType
	Owner *player.Player
}

func (p *Piece) String() string {
	return fmt.Sprintf("%s %s", p.Owner.Color, p.Type)
}
