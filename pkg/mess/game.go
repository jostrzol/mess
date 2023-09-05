package mess

type Game struct {
	*State
	controller Controller
	interactor Interactor
}

func NewGame(state *State, controller Controller, interactor Interactor) *Game {
	return &Game{
		State:      state,
		controller: controller,
		interactor: interactor,
	}
}

func (g *Game) PickWinner() (bool, *Player) {
	return g.controller.PickWinner(g.State)
}

func (g *Game) Choose(options []string) int {
	return g.interactor.ChooseOption(options)
}

func (g *Game) Move() error {
	move, err := g.interactor.ChooseMove(g.State, g.ValidMoves())
	if err != nil {
		return err
	}
	return move.Perform()
}

func (g *Game) Run() (*Player, error) {
	var winner *Player
	isFinished := false
	for !isFinished {
		g.interactor.PreTurn(g.State)

		err := g.controller.Turn(g.State)
		if err != nil {
			return nil, err
		}

		g.EndTurn()
		isFinished, winner = g.PickWinner()
	}
	return winner, nil
}

type Controller interface {
	PickWinner(state *State) (bool, *Player)
	Turn(state *State) error
}

type Interactor interface {
	PreTurn(state *State)
	ChooseOption(options []string) int
	ChooseMove(state *State, validMoves []Move) (*Move, error)
}
