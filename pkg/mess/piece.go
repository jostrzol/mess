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

func (p *Piece) GetCapturedBy(player *Player) error {
	if p.IsOnBoard() {
		return p.board.CaptureAt(p.square, player)
	}
	return nil
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
	for _, generated := range t.moveGenerators.Generate(piece) {
		move := Move{
			Piece:  piece,
			From:   piece.Square(),
			To:     generated.destination,
			Action: generated.action,
		}
		result = append(result, move)
	}
	return result
}

type MoveAction = func(*Piece, board.Square, board.Square)
type MoveGenerator func(*Piece) ([]board.Square, MoveAction)

type chainMoveGenerators []MoveGenerator
type moveGeneratorResult struct {
	destination board.Square
	action      MoveAction
}

func (g chainMoveGenerators) Generate(piece *Piece) []moveGeneratorResult {
	resultSet := make(map[brd.Square][]MoveAction, 0)
	for _, generator := range g {
		destinations, action := generator(piece)
		for _, destination := range destinations {
			if action != nil {
				resultSet[destination] = append(resultSet[destination], action)
			} else if _, present := resultSet[destination]; !present {
				resultSet[destination] = make([]MoveAction, 0)
			}
		}
	}
	result := make([]moveGeneratorResult, 0, len(resultSet))
	for destination, actions := range resultSet {
		// copy is required, because else the action closure
		// would always take actions from the last resultSet entry
		// (the 'action' reference changes as the loop iterates)
		actionsCopy := actions
		result = append(result, moveGeneratorResult{
			destination: destination,
			action: func(p *Piece, from, to brd.Square) {
				for _, action := range actionsCopy {
					action(piece, from, to)
				}
			},
		})
	}
	return result
}

type Move struct {
	Piece  *Piece
	From   brd.Square
	To     brd.Square
	Action MoveAction
}

func (m *Move) PerformWithoutAction() error {
	return m.Piece.MoveTo(m.To)
}

func (m *Move) Perform() error {
	err := m.Piece.MoveTo(m.To)
	if err != nil {
		return err
	}
	m.Action(m.Piece, m.From, m.To)
	return nil
}

func (m *Move) String() string {
	return fmt.Sprintf("%v: %v->%v", m.Piece, m.From, m.To)
}
