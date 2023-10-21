package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"sync"

	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/mess"
)

type interactor struct {
	scanner        *bufio.Scanner
	game           *mess.Game
	statePrinted   chan (struct{})
	movesGenerated chan (struct{})
	moveChosen     chan (*mess.Move)
	mutex          sync.Mutex
}

func newInteractor(game *mess.Game, scanner *bufio.Scanner) *interactor {
	return &interactor{
		scanner:        scanner,
		game:           game,
		statePrinted:   make(chan (struct{})),
		movesGenerated: make(chan (struct{})),
		moveChosen:     make(chan (*mess.Move)),
		mutex:          sync.Mutex{},
	}
}

func (t *interactor) Run() (*mess.Player, error) {
	var winner *mess.Player
	isFinished := false
	go t.chooseMove()
	t.printState()
	t.preloadMoves()
	for !isFinished {
		move := <-t.moveChosen
		if move == nil {
			t.closeChannels()
			return nil, ErrEOT
		}

		err := move.Perform()
		if err != nil {
			t.closeChannels()
			return nil, fmt.Errorf("performing move: %v", err)
		}

		t.game.EndTurn()
		t.printState()

		func() {
			t.mutex.Lock()
			defer t.mutex.Unlock()

			isFinished, winner = t.game.PickWinner()
		}()

		if !isFinished {
			go t.preloadMoves()
		}
	}
	t.closeChannels()
	return winner, nil
}

func (t *interactor) closeChannels() {
	close(t.statePrinted)
	close(t.movesGenerated)
	close(t.moveChosen)
}

func (t *interactor) printState() {
	fmt.Println(t.game.PrettyString())
	if len(t.statePrinted) == 0 {
		t.statePrinted <- struct{}{}
	}
}

func (t *interactor) preloadMoves() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.game.ValidMoves()
	if len(t.movesGenerated) == 0 {
		t.movesGenerated <- struct{}{}
	}
}

func (t *interactor) chooseMove() {
	moreWork := true
	for moreWork {
		_, moreWork = <-t.statePrinted
		move, err := t.tryChooseMove()
		if errors.Is(err, ErrCancel) {
			fmt.Println("<cancel>")

			go t.printState()
			go t.preloadMoves()
		} else if errors.Is(err, ErrEOT) {
			fmt.Println("<end of text>")
			moreWork = false
			t.moveChosen <- nil
		} else if err != nil {
			fmt.Printf("Error: %v!\n", err)
			fmt.Printf("Press enter to continue...")
			t.scanner.Scan()

			go t.printState()
			go t.preloadMoves()
		} else {
			t.moveChosen <- move
		}
	}
}

func (t *interactor) tryChooseMove() (*mess.Move, error) {
	println("Choose a square with your piece")
	print("> ")
	piece, err := t.chooseOwnPiece()
	if err != nil {
		return nil, err
	}

	<-t.movesGenerated
	validForPiece := make(map[board.Square]mess.MoveGroup, 0)
	for _, moveGroup := range t.game.ValidMoves() {
		if moveGroup.Piece == piece {
			validForPiece[moveGroup.To] = moveGroup
		}
	}

	if len(validForPiece) == 0 {
		return nil, fmt.Errorf("no valid moves for this piece")
	}
	println("Valid destinations:")
	for _, validMove := range validForPiece {
		fmt.Printf("-> %v\n", &validMove.To)
	}

	println("Choose a destination square")
	print("> ")
	dst, err := t.chooseSquare()
	if err != nil {
		return nil, err
	}
	moveGroup, ok := validForPiece[dst]
	if !ok {
		return nil, fmt.Errorf("invalid move")
	}

	for i := 0; i < moveGroup.ChoicesNumber(); i++ {
		moveGroup, err = t.chooseOption(moveGroup, i)
		if err != nil {
			return nil, err
		}
	}

	move := moveGroup.Single()
	return &move, nil
}

func (t *interactor) chooseOwnPiece() (*mess.Piece, error) {
	square, err := t.chooseSquare()
	if err != nil {
		return nil, err
	}

	piece, _ := t.game.Board().At(square)
	if piece == nil {
		return nil, fmt.Errorf("square is empty")
	} else if piece.Owner() != t.game.CurrentPlayer() {
		return nil, fmt.Errorf("piece belongs to the opponent")
	}

	return piece, nil
}

func (t *interactor) chooseSquare() (board.Square, error) {
	squareStr, err := t.scan()
	if err != nil {
		return board.Square{}, err
	}
	square, err := board.NewSquare(squareStr)
	if err != nil {
		return board.Square{}, err
	} else if !t.game.Board().Contains(square) {
		return board.Square{}, fmt.Errorf("Square not on board")
	}
	return square, nil
}

func (t *interactor) chooseOption(moveGroup mess.MoveGroup, choiceIdx int) (mess.MoveGroup, error) {
	optionStrings := moveGroup.UniqueOptionStrings(choiceIdx)

	println("Choose option:")
	for i, optionString := range optionStrings {
		fmt.Printf("%v. %v\n", i+1, optionString)
	}
	print("> ")
	choiceStr, err := t.scan()
	if err != nil {
		return moveGroup, err
	}

	var i int
	_, err = fmt.Sscanf(choiceStr, "%d", &i)
	if err != nil {
		return moveGroup, err
	}

	if i < 1 || i > len(optionStrings) {
		return moveGroup, fmt.Errorf("invalid option")
	}
	optionString := optionStrings[i-1]

	moveGroup = moveGroup.FilterMovesByOptionString(choiceIdx, optionString)
	return moveGroup, nil
}

func (t *interactor) scan() (string, error) {
	if !t.scanner.Scan() {
		if t.scanner.Err() == nil {
			return "", ErrEOT
		}
		return "", t.scanner.Err()
	}
	text := t.scanner.Text()
	if text == "" {
		return "", ErrCancel
	}
	return text, nil
}

var ErrCancel = fmt.Errorf("cancel")
var ErrEOT = fmt.Errorf("EOT")
