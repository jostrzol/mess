package color

//go:generate enumer -type=Color -transform=snake
type Color int

const (
	White Color = iota
	Black
)
