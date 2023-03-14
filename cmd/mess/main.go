package main

import (
	"flag"
	"fmt"
	"log"
	"sort"

	"github.com/jostrzol/mess/config"
	"github.com/jostrzol/mess/pkg/board"
	brd "github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/mess"
)

func chooseSquare(board *mess.PieceBoard) brd.Square {
	var square brd.Square
	var squareStr string
	var err error
	for ok := false; !ok; {
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
		} else {
			ok = true
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

func chooseMove(board *mess.PieceBoard, moves []mess.Move) brd.Square {
	var move brd.Square
	for ok := false; !ok; {
		move = chooseSquare(board)
		if !contains(moves, move) {
			println("Not an allowed move")
			println("Try again")
		} else {
			ok = true
		}
	}
	return move
}

func contains(moves []mess.Move, destination board.Square) bool {
	for _, move := range moves {
		if move.To == destination {
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

	var winner *mess.Player
	isFinished := false
	for !isFinished {
		println("Board: (uppercase - white, lowercase - black)")
		println(state.PrettyString())
		println("Choose square with a piece")
		piece := choosePiece(state.Board())
		if piece.Owner() != state.CurrentPlayer() {
			println("You must choose your piece")
			continue
		}

		moves := piece.ValidMoves()
		if len(moves) == 0 {
			println("No moves for this piece")
			continue
		}

		sort.Slice(moves, func(i, j int) bool {
			iSq := moves[i].To
			jSq := moves[j].To
			if iSq.Rank == jSq.Rank {
				return iSq.File < jSq.File
			}
			return iSq.Rank < jSq.Rank
		})

		print("Possible moves: ")
		for _, move := range moves {
			fmt.Printf("%v ", &move.To)
		}
		println()
		move := chooseMove(state.Board(), moves)

		err = piece.MoveTo(move)
		if err != nil {
			log.Fatal(err)
		}

		state.EndTurn()
		isFinished, winner = controller.PickWinner(state)
	}

	if winner == nil {
		fmt.Printf("Draw!\n")
	} else {
		fmt.Printf("Winner is %v!\n", winner)
	}
}
