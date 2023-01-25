board {
  height = 10
  width = 10
}

pieces {
  piece "king" {}
  piece "queen" {}
  piece "tower" {}
  piece "knight" {}
  piece "bishop" {}
  piece "pawn" {}
}

initial_state {
  piece_placement "white" {
    A1 = king
  }
  piece_placement "black" {
    A2 = king
    B3 = queen
  }
}
