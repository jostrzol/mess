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

func (g *Game) TurnOptions() ([][]Option, error) {
	choices, err := g.controller.TurnChoices(g.State)
	if err != nil {
		return nil, err
	}
	return choicesToOptionSets(choices), nil
}

func (g *Game) Turn(options []Option) error {
	return g.controller.Turn(g.State, options)
}

type Controller interface {
	PickWinner(state *State) (bool, *Player)
	TurnChoices(state *State) ([]Choice, error)
	Turn(state *State, options []Option) error
}
