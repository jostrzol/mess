package mess

import (
	"fmt"

	brd "github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/event"
)

type PieceType struct {
	name             string
	motionGenerators MotionGenerators
}

func NewPieceType(name string) *PieceType {
	return &PieceType{
		name:             name,
		motionGenerators: make(MotionGenerators, 0),
	}
}

func (t *PieceType) Name() string {
	return t.name
}

func (t *PieceType) String() string {
	return t.Name()
}

func (t *PieceType) AddMotionGenerator(generator MotionGenerator) {
	t.motionGenerators = append(t.motionGenerators, generator)
}

func (t *PieceType) generateMotions(piece *Piece) []brd.Square {
	return t.motionGenerators.GenerateMotions(piece)
}

type Piece struct {
	ty     *PieceType
	owner  *Player
	board  *PieceBoard
	square brd.Square
}

func NewPiece(pieceType *PieceType, owner *Player) *Piece {
	return &Piece{
		ty:    pieceType,
		owner: owner,
	}
}

func (p *Piece) String() string {
	return fmt.Sprintf("%s %s", p.Color(), p.ty)
}

func (p *Piece) Type() *PieceType {
	return p.ty
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

func (p *Piece) Square() *brd.Square {
	if p.IsOnBoard() {
		duplicate := p.square
		return &duplicate
	}
	return nil
}

func (p *Piece) IsOnBoard() bool {
	return p.board != nil
}

func (p *Piece) PlaceOn(board *PieceBoard, square *brd.Square) error {
	return board.Place(p, square)
}

func (p *Piece) RemoveFromBoard() {
	if p.IsOnBoard() {
		p.board.RemoveAt(&p.square)
	}
}

func (p *Piece) GenerateMotions() []brd.Square {
	return p.ty.generateMotions(p)
}

func (p *Piece) MoveTo(square *brd.Square) error {
	return p.board.Move(p, square)
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
			p.square = e.ToSquare
		}
	case PieceRemoved:
		if e.Piece == p {
			p.board = nil
		}
	}
}