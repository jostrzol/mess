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
    C5 = "queen"
  }
  pieces "black" {
    A2 = "king"
    B3 = "queen"
    B5 = "pawn"
  }
}

variable "piece_points" {
  value = {
    king = 1000
    queen = 9
    rook = 5
    knight = 3
    bishop = 3
    pawn = 1
  }
}

function "calc_player_points" {
  params = [player]
  result = sum([for i, piece in player.pieces: piece_points[piece.type]]...)
}

function "calc_points_per_player" {
  params = [players]
  result = {for i, player in players: calc_player_points(player) => player...}
}

function "best_players" {
  params = [points_per_player]
  result = points_per_player[max(keys(points_per_player)...)]
}

function "pick_winner_or_draw" {
  params = [best_players]
  result = length(best_players) == 1 ? best_players[0] : null
}

function "decide_winner" {
  params = [game]
  result = pick_winner_or_draw(best_players(calc_points_per_player(game.players)))
}