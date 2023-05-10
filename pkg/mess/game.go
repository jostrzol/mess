package mess

type Game struct {
	*State
	controller Controller
}

func NewGame(state *State, controller Controller) *Game {
	return &Game{
		State:      state,
		controller: controller,
	}
}

func (g *Game) PickWinner() (bool, *Player) {
	return g.controller.PickWinner(g.State)
}

func (g *Game) Choose(options []string) int {
	return g.controller.Choose(options)
}

type Controller interface {
	PickWinner(state *State) (bool, *Player)
	Choose(options []string) int
}
