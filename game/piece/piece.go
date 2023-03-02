package piece

import (
	"fmt"

	brd "github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/player"
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
	Owner  *player.Player
	Board  Board
	Square *brd.Square
}

func NewPiece(pieceType *Type, owner *player.Player) *Piece {
	return &Piece{
		Type:  pieceType,
		Owner: owner,
	}
}

func (p *Piece) String() string {
	return fmt.Sprintf("%s %s", p.Owner.Color, p.Type)
}

func (p *Piece) PlaceOn(board Board, square brd.Square) error {
	err := board.Place(p, &square)
	if err != nil {
		return err
	}
	p.Board = board
	p.Square = &square
	return nil
}

func (p *Piece) GenerateMotions() []brd.Square {
	return p.Type.generateMotions(p)
}

func (p *Piece) MoveTo(square brd.Square) error {
	// TODO: add capturing on destination square
	if p.Board == nil {
		return fmt.Errorf("piece not on board")
	}
	err := p.Board.Place(p, &square)
	if err != nil {
		return err
	}
	return nil
}
