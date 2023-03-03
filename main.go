package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jostrzol/mess/config"
	brd "github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/piece"
)

func chooseSquare(board piece.Board) *brd.Square {
	var square *brd.Square
	var squareStr string
	var err error
	for square == nil || err != nil {
		print("> ")
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

func choosePiece(board piece.Board) *piece.Piece {
	var piece *piece.Piece
	var err error
	for piece == nil || err != nil {
		square := chooseSquare(board)
		piece, err = board.At(square)
		if err != nil {
			fmt.Printf("%v\n", err)
			println("Try again")
		} else if piece == nil {
			println("No piece there")
			println("Try again")
		}
	}
	return piece
}

func chooseMove(board piece.Board, moves []brd.Square) *brd.Square {
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

	print(state.String())
	println("Choose square with a piece")
	piece := choosePiece(state.Board)

	motions := piece.GenerateMotions()
	if len(motions) == 0 {
		log.Fatal("no motions for this piece")
	}

	fmt.Printf("Possible moves: %v", motions)
	move := chooseMove(state.Board, motions)

	err = state.Move(piece, move)
	if err != nil {
		log.Fatal(err)
	}

	winner := controller.DecideWinner(state)
	if winner == nil {
		fmt.Printf("Draw!\n")
	} else {
		fmt.Printf("Winner is %v!\n", winner)
	}
}
