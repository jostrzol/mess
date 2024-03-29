package mess

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/board"
	brd "github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/event"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Player struct {
	color            color.Color
	pieces           map[*Piece]struct{}
	captures         map[*Piece]struct{}
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
		player.captures = make(map[*Piece]struct{})
		board.Observe(player)
	}
	return players
}

func (p *Player) Color() color.Color {
	return p.color
}

func (p *Player) Pieces() []*Piece {
	return maps.Keys(p.pieces)
}

func (p *Player) Captures() []*Piece {
	return maps.Keys(p.captures)
}

func (p *Player) CapturesCountByType() map[*PieceType]int {
	result := make(map[*PieceType]int, 0)
	for piece := range p.captures {
		result[piece.ty]++
	}
	return result
}

func (p *Player) ConvertAndReleasePiece(pieceType *PieceType, board *PieceBoard, square brd.Square) error {
	captures := p.Captures()
	i := slices.IndexFunc(captures, func(piece *Piece) bool { return piece.Type() == pieceType })
	if i == -1 {
		return fmt.Errorf("player %v does not have a capture of type %v", p, pieceType)
	}

	piece := captures[i]
	oldOwner := piece.owner
	piece.owner = p

	err := board.Place(piece, square)
	if err != nil {
		piece.owner = oldOwner
		return fmt.Errorf("placing the released piece at destination square: %w", err)
	}

	return nil
}

func (p *Player) ForwardDirection() board.Offset {
	return p.forwardDirection
}

func (p *Player) String() string {
	return p.color.String()
}

func (p *Player) Moves() []*MoveGroup {
	result := make([]*MoveGroup, 0)
	for piece := range p.pieces {
		result = append(result, piece.Moves()...)
	}
	return result
}

func (p *Player) AttackedSquares() []board.Square {
	resultSet := make(map[board.Square]struct{})
	for piece := range p.pieces {
		for _, move := range piece.Moves() {
			resultSet[move.To] = struct{}{}
		}
	}
	result := make([]board.Square, 0, len(resultSet))
	for square := range resultSet {
		result = append(result, square)
	}
	return result
}

func (p *Player) Handle(event event.Event) {
	switch e := event.(type) {
	case PiecePlaced:
		if e.Piece.Owner() == p {
			p.pieces[e.Piece] = struct{}{}
		}
		delete(p.captures, e.Piece)
	case PieceCaptured:
		if e.CapturedBy == p {
			p.captures[e.Piece] = struct{}{}
		}
	case PieceRemoved:
		if e.Piece.Owner() == p {
			delete(p.pieces, e.Piece)
		}
	}
}
