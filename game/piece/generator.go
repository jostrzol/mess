package piece

import brd "github.com/jostrzol/mess/game/board"

type MotionGenerator interface {
	GenerateMotions(piece *Piece) []brd.Square
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
