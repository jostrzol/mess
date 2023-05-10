package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jostrzol/mess/config"
	"github.com/jostrzol/mess/pkg/board"
	brd "github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/mess"
)

type terminalInteractor struct{}

func (i terminalInteractor) Choose(options []string) int {
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

func chooseMove(game *mess.Game) *mess.Move {
	validMoves := game.ValidMoves()
	for {
		var srcStr string
		println("Choose a square with your piece")
		print("> ")
		fmt.Scanf("%s", &srcStr)

		src, err := brd.NewSquare(srcStr)
		if err != nil {
			fmt.Printf("%v\n", err)
			println("Try again")
			continue
		} else if !game.Board().Contains(src) {
			println("Square not on board")
			println("Try again")
			continue
		}

		piece, _ := game.Board().At(src)
		if piece == nil {
			println("That square is empty!")
			continue
		} else if piece.Owner() != game.CurrentPlayer() {
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

		var dstStr string
		println("Choose a destination square")
		print("> ")
		fmt.Scanf("%s", &dstStr)

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

		return &move
	}
}

func main() {
	var configFilename = flag.String("rules", "./rules.hcl", "path to a rules config file")
	flag.Parse()

	game, err := config.DecodeConfig(*configFilename, terminalInteractor{})
	if err != nil {
		log.Fatalf("loading game rules: %s", err)
	}

	var winner *mess.Player
	isFinished := false
	for !isFinished {
		// generate moves first so that debug logs print before the board does
		// (the moves are cached anyway, so this computation won't get wasted)
		game.ValidMoves()

		println("Board: (uppercase - white, lowercase - black)")
		println(game.PrettyString())

		move := chooseMove(game)

		err = move.Perform()
		if err != nil {
			log.Fatal(err)
		}

		game.EndTurn()
		isFinished, winner = game.PickWinner()
	}

	if winner == nil {
		fmt.Printf("Draw!\n")
	} else {
		fmt.Printf("Winner is %v!\n", winner)
	}
}
