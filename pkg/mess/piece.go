package mess

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/board"
	brd "github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/event"
)

type Piece struct {
	ty     *PieceType
	owner  *Player
	board  *PieceBoard
	square brd.Square
	moves  []Move
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

func (p *Piece) Square() brd.Square {
	return p.square
}

func (p *Piece) IsOnBoard() bool {
	return p.board != nil
}

func (p *Piece) PlaceOn(board *PieceBoard, square brd.Square) error {
	return board.Place(p, square)
}

func (p *Piece) MoveTo(square brd.Square) error {
	return p.board.Move(p, square)
}

func (p *Piece) RemoveFromBoard() {
	if p.IsOnBoard() {
		p.board.RemoveAt(p.square)
	}
}

func (p *Piece) Moves() []Move {
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

type PieceType struct {
	name           string
	moveGenerators chainMoveGenerators
}

func NewPieceType(name string) *PieceType {
	return &PieceType{
		name:           name,
		moveGenerators: make(chainMoveGenerators, 0),
	}
}

func (t *PieceType) Name() string {
	return t.name
}

func (t *PieceType) String() string {
	return t.Name()
}

func (t *PieceType) AddMoveGenerator(generator MoveGenerator) {
	t.moveGenerators = append(t.moveGenerators, generator)
}

func (t *PieceType) moves(piece *Piece) []Move {
	result := make([]Move, 0)
	for _, destination := range t.moveGenerators.Generate(piece) {
		move := Move{
			Piece: piece,
			From:  piece.Square(),
			To:    destination,
		}
		result = append(result, move)
	}
	return result
}

type MoveGenerator func(*Piece) []board.Square
type chainMoveGenerators []MoveGenerator

func (g chainMoveGenerators) Generate(piece *Piece) []brd.Square {
	destinationSet := make(map[brd.Square]bool, 0)
	for _, generator := range g {
		newDestinations := generator(piece)
		for _, destination := range newDestinations {
			destinationSet[destination] = true
		}
	}
	destinations := make([]brd.Square, 0, len(destinationSet))
	for s := range destinationSet {
		destinations = append(destinations, s)
	}
	return destinations
}

type Move struct {
	Piece *Piece
	From  brd.Square
	To    brd.Square
}

func (m *Move) Perform() error {
	return m.Piece.MoveTo(m.To)
}

func (m *Move) String() string {
	return fmt.Sprintf("%v: %v->%v", m.Piece, m.From, m.To)
}
