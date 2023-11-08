package room

import "github.com/jostrzol/mess/pkg/mess"

type State struct {
	Board      *mess.PieceBoard
	OptionTree *mess.OptionNode
}
