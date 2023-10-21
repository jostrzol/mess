package cmd

import (
	"bufio"
	"io"

	"github.com/jostrzol/mess/pkg/mess"
)

func Run(game *mess.Game, in io.Reader, out io.Writer) (*mess.Player, error) {
	scanner := bufio.NewScanner(in)
	// TODO: handle out
	i := newInteractor(game, scanner)
	return i.Run()
}
