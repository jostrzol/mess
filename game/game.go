package game

type Game struct {
	Board            Board
	Players          []*Player
	DecideWinnerFunc func(*Game) *Player
}

func (g *Game) DecideWinner() *Player {
	return g.DecideWinnerFunc(g)
}
