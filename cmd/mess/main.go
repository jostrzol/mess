package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/jostrzol/mess/pkg/board"
	brd "github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/jostrzol/mess/pkg/rules"
)

type terminalInteractor struct {
	scanner *bufio.Scanner
}

func (t *terminalInteractor) PreTurn(state *mess.State) {
	// generate moves first so that debug logs print before the board does
	// (the moves are cached anyway, so this computation won't get wasted)
	state.ValidMoves()

	println("Board: (uppercase - white, lowercase - black)")
	println(state.PrettyString())

}

func (t *terminalInteractor) ChooseOption(options []string) int {
	println("Choose option:")
	for i, option := range options {
		fmt.Printf("%d. %s\n", i+1, option)
	}
	var choice int
	for {
		_, err := fmt.Scanf("%d", &choice)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			println("Try again")
		} else if choice > len(options) || choice == 0 {
			fmt.Printf("incorect choice\n")
			println("Try again")
		} else {
			return choice - 1
		}
	}
}

func (t *terminalInteractor) ChooseMove(state *mess.State, validMoves []mess.Move) (*mess.Move, error) {
	for {
		println("Choose a square with your piece")
		print("> ")
		if !t.scanner.Scan() {
			return nil, t.scanner.Err()
		}
		srcStr := t.scanner.Text()
		if srcStr == "" {
			continue
		}

		src, err := brd.NewSquare(srcStr)
		if err != nil {
			fmt.Printf("%v\n", err)
			println("Try again")
			continue
		} else if !state.Board().Contains(src) {
			println("Square not on board")
			println("Try again")
			continue
		}

		piece, _ := state.Board().At(src)
		if piece == nil {
			println("That square is empty!")
			continue
		} else if piece.Owner() != state.CurrentPlayer() {
			println("That belongs to your opponent!")
			continue
		}

		validForPiece := make(map[board.Square]mess.Move, 0)
		for _, validMove := range validMoves {
			if validMove.Piece == piece {
				validForPiece[validMove.To] = validMove
			}
		}

		if len(validForPiece) == 0 {
			println("No valid moves for this piece!")
			continue
		} else {
			println("Valid destinations:")
			for _, validMove := range validForPiece {
				fmt.Printf("-> %v\n", &validMove.To)
			}
		}

		println("Choose a destination square")
		print("> ")
		if !t.scanner.Scan() {
			return nil, t.scanner.Err()
		}
		dstStr := t.scanner.Text()
		if dstStr == "" {
			continue
		}

		dst, err := brd.NewSquare(dstStr)
		if err != nil {
			fmt.Printf("%v\n", err)
			println("Try again")
		}

		move, ok := validForPiece[dst]
		if !ok {
			println("Invalid move!")
			continue
		}

		return &move, nil
	}
}

func cmdError(format string, a ...any) {
	format = fmt.Sprintf("error: %s\n", format)
	fmt.Printf(format, a...)
	flag.Usage()
	os.Exit(1)
}

func runError(format string, a ...any) {
	format = fmt.Sprintf("error: %s\n", format)
	fmt.Printf(format, a...)
	os.Exit(2)
}

func main() {
	var rulesFilename = flag.String("rules", "", "path to a rules file")
	flag.Parse()

	if *rulesFilename == "" {
		cmdError("no rules file")
	}

	scanner := bufio.NewScanner(os.Stdin)
	interactor := &terminalInteractor{scanner}
	game, err := rules.DecodeRules(*rulesFilename, interactor, true)
	if err != nil {
		runError("loading game rules: %s", err)
	}

	winner, err := game.Run()
	if err != nil {
		runError("running game: %s", err)
	}

	if winner == nil {
		fmt.Printf("Draw!\n")
	} else {
		fmt.Printf("Winner is %v!\n", winner)
	}
}
