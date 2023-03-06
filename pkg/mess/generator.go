package mess

import brd "github.com/jostrzol/mess/pkg/board"

type MotionGenerator interface {
	GenerateMotions(piece *Piece) []brd.Square
}

type FuncMotionGenerator func(piece *Piece) []brd.Square

func (g FuncMotionGenerator) GenerateMotions(piece *Piece) []brd.Square {
	return g(piece)
}

type MotionGenerators []MotionGenerator

func (g MotionGenerators) GenerateMotions(piece *Piece) []brd.Square {
	destinationSet := make(map[brd.Square]bool, 0)
	for _, generator := range g {
		newDestinations := generator.GenerateMotions(piece)
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