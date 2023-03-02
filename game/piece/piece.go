package piece

import (
	"fmt"
	"log"

	brd "github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/piece/color"
)

type Board = brd.Board[*Piece]

type Type struct {
	Name             string
	motionGenerators MotionGenerators
}

func NewType(name string) *Type {
	return &Type{
		Name:             name,
		motionGenerators: make(MotionGenerators, 0),
	}
}

func (t *Type) String() string {
	return t.Name
}

func (t *Type) AddMotionGenerator(generator MotionGenerator) {
	t.motionGenerators = append(t.motionGenerators, generator)
}

func (t *Type) generateMotions(piece *Piece) []brd.Square {
	return t.motionGenerators.GenerateMotions(piece)
}

type Piece struct {
	Type   *Type
	Owner  Owner
	Board  Board
	Square brd.Square
}

type Owner interface {
	Color() color.Color
}

func NewPiece(pieceType *Type, owner Owner) *Piece {
	return &Piece{
		Type:  pieceType,
		Owner: owner,
	}
}

func (p *Piece) String() string {
	return fmt.Sprintf("%s %s", p.Color(), p.Type)
}

func (p *Piece) Color() color.Color {
	return p.Owner.Color()
}

func (p *Piece) PlaceOn(board Board, square *brd.Square) error {
	old, err := board.Place(p, square)
	if err != nil {
		return err
	}
	if old != nil {
		old.Board = nil
		log.Printf("replacing %v with %v on %v", old, p, &square)
	}
	p.Board = board
	p.Square = *square
	return nil
}

func (p *Piece) GenerateMotions() []brd.Square {
	return p.Type.generateMotions(p)
}

func (p *Piece) MoveTo(square *brd.Square) (*Piece, error) {
	if p.Board == nil {
		return nil, fmt.Errorf("piece not on board")
	}

	_, err := p.Board.Place(nil, &p.Square)
	if err != nil {
		return nil, err
	}
	old, err := p.Board.Place(p, square)
	if err != nil {
		p.Board.Place(p, &p.Square)
		return nil, err
	}

	p.Square = *square
	if old != nil {
		old.Board = nil
	}
	return old, nil
}
