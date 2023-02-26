package piece

func Rook() *Type {
	return &Type{
		Name: "rook",
	}
}

func Knight() *Type {
	return &Type{
		Name: "knight",
	}
}

func Noones(pieceType *Type) *Piece {
	return &Piece{
		Type:  pieceType,
		Owner: nil,
	}
}
