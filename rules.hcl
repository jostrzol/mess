board {
  height = 8
  width = 8
}

piece_types {
  piece_type "king" {}
  piece_type "queen" {}
  piece_type "rook" {}
  piece_type "knight" {}
  piece_type "bishop" {}
  piece_type "pawn" {}
}

initial_state {
  pieces "white" {
    A1 = "king"
  }
  pieces "black" {
    A2 = "king"
    B3 = "queen"
  }
}

function "calculate_points" {
  params = [player]
  result = player.color == "white" ? 7 : 3
}

function "points_to_players" {
  params = [players]
  result = {for i, player in players: calculate_points(player) => player}
}

function "max_points" {
  params = [players]
  result = max([for i, player in players: calculate_points(player)]...)
}

function "decide_winner" {
  params = [game]
  result = points_to_players(game.players)[max_points(game.players)]
}