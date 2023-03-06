package mess

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/event"
)

type PieceBoard struct {
	event.Subject
	wrapped board.Board[*Piece]
}

func NewPieceBoard(width int, height int) (*PieceBoard, error) {
	board, err := board.NewBoard[*Piece](width, height)
	if err != nil {
		return nil, err
	}
	return &PieceBoard{
		Subject: event.NewSubject(),
		wrapped: board,
	}, nil
}

func (b *PieceBoard) String() string {
	return b.wrapped.String()
}

func (b *PieceBoard) At(square *board.Square) (*Piece, error) {
	return b.wrapped.At(square)
}

func (b *PieceBoard) Contains(square *board.Square) bool {
	return b.wrapped.Contains(square)
}

func (b *PieceBoard) AllPieces() []*Piece {
	return b.wrapped.AllItems()
}

func (b *PieceBoard) Place(piece *Piece, square *board.Square) error {
	if piece.IsOnBoard() {
		return fmt.Errorf("piece already on a board")
	}
	old, err := b.wrapped.At(square)
	if err != nil {
		return fmt.Errorf("getting piece at %v: %w", square, err)
	}
	if old != nil {
		return fmt.Errorf("placing %v on %v: already occupied by %v", piece, square, old)
	}

	_, err = b.wrapped.Place(piece, square)
	if err != nil {
		return err
	}
	b.Observe(piece)
	b.Notify(PiecePlaced{
		Piece:  piece,
		Board:  b,
		Square: *square,
	})
	return nil
}

type PiecePlaced struct {
	Piece  *Piece
	Board  *PieceBoard
	Square board.Square
}

func (b *PieceBoard) RemoveAt(square *board.Square) error {
	old, err := b.wrapped.Place(nil, square)
	if err != nil {
		return err
	}
	if old != nil {
		b.notifyRemoved(old)
	}
	return nil
}

func (b *PieceBoard) notifyRemoved(piece *Piece) {
	b.Notify(PieceRemoved{
		Piece:  piece,
		Square: *piece.Square(),
	})
	b.Unobserve(piece)
}

type PieceRemoved struct {
	Piece  *Piece
	Square board.Square
}

func (b *PieceBoard) Move(piece *Piece, square *board.Square) error {
	if piece.Board() != b {
		return fmt.Errorf("piece not on board")
	}

	_, err := b.wrapped.Place(nil, piece.Square())
	if err != nil {
		return err
	}
	old, err := b.wrapped.Place(piece, square)
	if err != nil {
		b.wrapped.Place(piece, piece.Square())
		return err
	}

	b.Notify(PieceMoved{
		Piece:      piece,
		FromSquare: *piece.Square(),
		ToSquare:   *square,
	})

	if old != nil {
		b.notifyRemoved(old)
		b.Notify(PieceCaptured{
			Piece:        old,
			CapturedFrom: old.Owner(),
			CapturedBy:   piece.Owner(),
		})
	}
	return nil
}

type PieceMoved struct {
	Piece      *Piece
	FromSquare board.Square
	ToSquare   board.Square
}
type PieceCaptured struct {
	Piece        *Piece
	CapturedBy   *Player
	CapturedFrom *Player
}
