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

func (g *Game) TurnOptions() (*OptionNode, error) {
	choice, err := g.controller.TurnChoice(g.State)
	if err != nil {
		return nil, err
	}
	return choice.GenerateOptions(), nil
}

func (g *Game) Turn(options []Option) error {
	err := g.controller.Turn(g.State, options)
	if err != nil {
		return err
	}
	g.State.EndTurn()
	return nil
}

type Controller interface {
	PickWinner(state *State) (bool, *Player)
	TurnChoice(state *State) (*Choice, error)
	Turn(state *State, options []Option) error
}
