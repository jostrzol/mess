board {
  height = 8
  width = 8
}

piece_types {
  piece_type king {}
  piece_type queen {}
  piece_type rook {}
  piece_type knight {}
  piece_type bishop {}
  piece_type pawn {}
}

initial_state {
  pieces white {
    A1 = "king"
  }
  pieces black {
    A2 = "king"
    B3 = "queen"
  }
}

function "decide_winner" {
  params = [game]
  result = game.players.white
}