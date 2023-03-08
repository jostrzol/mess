package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jostrzol/mess/config"
	brd "github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/mess"
)

func chooseSquare(board *mess.PieceBoard) *brd.Square {
	var square *brd.Square
	var squareStr string
	var err error
	for square == nil || err != nil {
		print("> ")
		// squareStr = "A1"
		fmt.Scan(&squareStr)
		square, err = brd.NewSquare(squareStr)
		if err != nil {
			fmt.Printf("%v\n", err)
			println("Try again")
		} else if !board.Contains(square) {
			println("Square not in board")
			println("Try again")
			square = nil
		}
	}
	return square
}

func choosePiece(board *mess.PieceBoard) *mess.Piece {
	var piece *mess.Piece
	var err error
	for piece == nil || err != nil {
		square := chooseSquare(board)
		piece, err = board.At(square)
		if err != nil {
			println(err)
			println("Try again")
		}
		if piece == nil {
			println("No piece there")
			println("Try again")
		}
	}
	return piece
}

func chooseMove(board *mess.PieceBoard, moves []brd.Square) *brd.Square {
	var move *brd.Square
	var err error
	for move == nil || err != nil {
		move = chooseSquare(board)
		if !contains(moves, *move) {
			println("Not an allowed move")
			println("Try again")
			move = nil
		}
	}
	return move
}

func contains[T comparable](slice []T, item T) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}

func main() {
	var configFilename = flag.String("rules", "./rules.hcl", "path to a rules config file")
	flag.Parse()

	state, controller, err := config.DecodeConfig(*configFilename)
	if err != nil {
		log.Fatalf("loading game rules: %s", err)
	}

	for {
		println("Board: (uppercase - white, lowercase - black)")
		println(state.PrettyString())
		println("Choose square with a piece")
		piece := choosePiece(state.Board())
		if piece.Owner() != state.CurrentPlayer() {
			println("You must choose your piece")
			continue
		}

		motions := piece.GenerateMotions()
		if len(motions) == 0 {
			println("No motions for this piece")
			continue
		}

		print("Possible moves: ")
		for _, motion := range motions {
			fmt.Printf("%v ", &motion)
		}
		println()
		move := chooseMove(state.Board(), motions)

		err = piece.MoveTo(move)
		if err != nil {
			log.Fatal(err)
		}

		state.EndTurn()
	}

	winner := controller.DecideWinner(state)
	if winner == nil {
		fmt.Printf("Draw!\n")
	} else {
		fmt.Printf("Winner is %v!\n", winner)
	}
}
