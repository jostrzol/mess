board {
  height = 8
  width = 8
}

pieces {
  piece king {}
  piece queen {}
  piece rook {}
  piece knight {}
  piece bishop {}
  piece pawn {}
}

initial_state {
  piece_placements white {
    A1 = king
  }
  piece_placements black {
    A2 = king
    B3 = queen
  }
}

function "decide_winner" {
  params = [game]
  result = game.players[0]
}