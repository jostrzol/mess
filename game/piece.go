package game

type PieceType struct {
	Name string
}

type Piece struct {
	Type  PieceType
	Owner Player
}
