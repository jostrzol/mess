package mess

import (
	"fmt"

	brd "github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/event"
)

type Piece struct {
	ty     *PieceType
	owner  *Player
	board  *PieceBoard
	square brd.Square
	moves  []*MoveGroup
}

func NewPiece(pieceType *PieceType, owner *Player) *Piece {
	return &Piece{
		ty:    pieceType,
		owner: owner,
	}
}

func (p *Piece) String() string {
	var colorStr string
	if p.owner != nil {
		colorStr = p.Color().String()
	} else {
		colorStr = "noone's"
	}
	return fmt.Sprintf("%s %s", colorStr, p.ty)
}

func (p *Piece) Type() *PieceType {
	return p.ty
}

func (p *Piece) Presentation() Presentation {
	if p.owner == nil {
		return Presentation{Symbol: rune('?')}
	}
	return p.ty.Presentation(p.Color())
}

func (p *Piece) Owner() *Player {
	return p.owner
}

func (p *Piece) Color() color.Color {
	return p.owner.Color()
}

func (p *Piece) Board() *PieceBoard {
	return p.board
}

func (p *Piece) Square() brd.Square {
	return p.square
}

func (p *Piece) IsOnBoard() bool {
	return p.board != nil
}

func (p *Piece) PlaceOn(board *PieceBoard, square brd.Square) error {
	return board.Place(p, square)
}

func (p *Piece) Remove() error {
	return p.board.RemoveAt(p.square)
}

func (p *Piece) MoveTo(square brd.Square) error {
	return p.board.Move(p, square)
}

func (p *Piece) GetCapturedBy(player *Player) error {
	if p.IsOnBoard() {
		return p.board.CaptureAt(p.square, player)
	}
	return nil
}

func (p *Piece) Moves() []*MoveGroup {
	if p.moves == nil {
		p.generateMoves()
	}
	return p.moves
}

func (p *Piece) generateMoves() {
	p.moves = p.ty.moves(p)
}

func (p *Piece) Handle(event event.Event) {
	switch e := event.(type) {
	case PiecePlaced:
		if e.Piece == p {
			p.board = e.Board
			p.square = e.Square
		}
	case PieceMoved:
		if e.Piece == p {
			p.square = e.To
		}
	case PieceRemoved:
		if e.Piece == p {
			p.board = nil
		}
	}
	p.moves = nil
}

func (p *Piece) Clone() *Piece {
	return NewPiece(p.ty, p.owner)
}
