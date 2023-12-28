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

func (g *Game) TurnOptions() (*OptionNode, error) {
	choice, err := g.controller.TurnChoice(g.State)
	if err != nil {
		return nil, err
	}
	optionTree := choice.GenerateOptions()
	return optionTree, nil
}

func (g *Game) PlayTurn(options []Option) error {
	err := g.controller.Turn(g.State, options)
	if err != nil {
		return err
	}
	g.State.EndTurn()
	return nil
}

func (g *Game) Resolution() Resolution {
	return g.controller.Resolution(g.State)
}

type Controller interface {
	TurnChoice(state *State) (*Choice, error)
	Turn(state *State, options []Option) error
	Resolution(state *State) Resolution
}

type Resolution struct {
	DidEnd bool
	Winner *Player
}
