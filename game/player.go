package game

type Color int

const (
	White Color = iota
	Black
)

type Player struct {
	Color Color
}
