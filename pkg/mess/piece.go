package mess

import (
	"fmt"
	"log"

	brd "github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/color"
)

type PieceBoard = brd.Board[*Piece]

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

type Owner interface {
	Color() color.Color
}

type Piece struct {
	ty     *PieceType
	owner  Owner
	board  brd.Board[*Piece]
	square brd.Square
}

func NewPiece(pieceType *PieceType, owner Owner) *Piece {
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

func (p *Piece) Owner() Owner {
	return p.owner
}

func (p *Piece) Color() color.Color {
	return p.owner.Color()
}

func (p *Piece) Board() PieceBoard {
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

func (p *Piece) PlaceOn(board PieceBoard, square *brd.Square) error {
	old, err := board.Place(p, square)
	if err != nil {
		return err
	}
	if old != nil {
		old.board = nil
		log.Printf("replacing %v with %v on %v", old, p, &square)
	}
	p.board = board
	p.square = *square
	return nil
}

func (p *Piece) RemoveFromBoard() {
	if p.board != nil {
		p.board.Place(nil, &p.square)
		p.board = nil
	}
}

func (p *Piece) GenerateMotions() []brd.Square {
	return p.ty.generateMotions(p)
}

func (p *Piece) MoveTo(square *brd.Square) (*Piece, error) {
	if p.board == nil {
		return nil, fmt.Errorf("piece not on board")
	}

	_, err := p.board.Place(nil, &p.square)
	if err != nil {
		return nil, err
	}
	old, err := p.board.Place(p, square)
	if err != nil {
		p.board.Place(p, &p.square)
		return nil, err
	}

	p.square = *square
	if old != nil {
		old.board = nil
	}
	return old, nil
}
