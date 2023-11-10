package room

import "github.com/jostrzol/mess/pkg/mess"

type State struct {
	State      *mess.State
	TurnNumber int
	Board      *mess.PieceBoard
	OptionTree *mess.OptionNode
}
