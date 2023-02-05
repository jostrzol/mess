board {
  height = 8
  width = 8
}

// ===== HELPER FUNCTIONS ======================================

// Checks if square is occupied by a piece of a given player
function "is_square_owned_by" {
  params = [square, player]
  result = square.piece != null && square.piece.owner == player
}

// ===== MOVE FUNCTIONS ========================================
// They receive 2 parameters:
//  * square - the current square
//  * piece - the current piece
// and generate all the squares that the given piece can move to from the given square.

// Generates moves to the straight neighbours (top, right, bottom, left)
// of the current square, given that they are not occupied by the player
// owning the current piece.
composite_function "move_neighbours_straight" {
  params = [square, piece]
  result = {
    dposes = [[0, 1], [1, 0], [0, -1], [-1, 0]]
    neighbours_straight = [get_square_relative(square, dpos) for dpos in dposes]
    return = [neighbour for neighbour in neighbours_straight if !is_square_owned_by(square, piece.owner)]
  }
}

// Generates a move one square forwards, given that the destination square
// is not occupied by any piece.
composite_function "move_forward_straight" {
  params = [square, piece]
  result = {
    dest = get_square_relative(square, piece.owner.forward_direction)
    return = dest != null && dest.piece == null ? [dest] : []
  }
}

// Generates a move two square forwards, given that the destination square
// is not occupied by any piece and that the piece has not moved yet before.
composite_function "move_forward_straight_double" {
  params = [square, piece]
  result = {
    dpos = [dcoord * 2 for dcoord in piece.owner.forward_direction]
    dest = get_square_relative(square, dpos)
    return = dest != null && dest.piece == null && !has_ever_moved(piece) ? [dest] : []
  }
}

// Generates 2 moves: one square forwards and to either side, given that the
// destination squares are not occupied by the player owning the current piece.
composite_function "move_forward_diagonal" {
  params = [square, piece]
  result = {
    forward_y = piece.owner.forward_direction[1]
    forward_left = [-1, forward_y]
    forward_right = [1, forward_y]
    dest1 = get_square_relative(square, forward_left)
    dest2 = get_square_relative(square, forward_right)
    dests = [dest1, dest2]
    return = [dest for dest in dests if dest != null && !is_square_owned_by(piece.owner)]
  }
}

// Generates a maximum of 8 moves, meeting criteria:
//   * first go 2 to any side,
//   * then go 1 to any side, but the direction is perpendicular to the one of previous move.
// If the destination square is occupied by the player owing the current piece, it is discarded.
composite_function "move_hook" {
  params = [square, piece]
  result = {
    dposes = [[2, 1], [2, -1], [-2, 1], [-2, -1], [1, 2], [-1, 2], [1, -2], [-1, -2]]
    dests = [get_square_relative(square, dpos) for dpos in dposes]
    return = [dest for dest in dests if dest != null && !is_square_owned_by(piece.owner)]
  }
}

// Generates moves from current position (param 'square') in the given direction (param 'dpos' in
// form [dx, dy]) until end of board or a piece is encountered. If said piece has the same owner as
// the one in param 'piece', the last square is excluded from the generated square, else it is
// included.
composite_function "move_line" {
  params = [square, piece, dpos]
  result = {
    next = get_square_relative(square, dpos)
    return = next == null ? [] : next.piece == null ? [next, moves_line(next, piece, dpos)...] : is_square_owned_by(next, piece.owner) ? [] : [next]
  }
}

composite_function "move_line_diagonal" {
  params = [square, piece]
  result = {
    dposes = [[-1, 1], [1, 1], [1, -1], [-1, -1]]
    return = [moves_line(square, piece, dpos)... for dpos in dposes]
  }
}

composite_function "move_line_straight" {
  params = [square, piece]
  result = {
    dposes = [[0, 1], [1, 0], [0, -1], [-1, 0]]
    return = [moves_line(square, piece, dpos)... for dpos in dposes]
  }
}

piece_types {
  piece_type "king" {
    // TODO: castling
    // TODO: block moving to dangerous squares
    move {
      generator = "move_neighbours_straight"
    }
  }

  piece_type "queen" {
    move {
      generator = "move_line_diagonal"
    }
    move {
      generator = "move_line_straight"
    }
  }

  piece_type "rook" {
    move {
      generator = "move_line_straight"
    }
  }

  piece_type "knight" {
    move {
      generator = "move_hook"
    }
  }

  piece_type "bishop" {
    move {
      generator = "move_line_diagonal"
    }
  }

  piece_type "pawn" {
    // TODO: en passant
    move {
      generator = "move_forward_straight"
    }
    move {
      generator = "move_forward_straight_double"
    }
    move {
      generator = "move_forward_diagonal"
    }
  }
}

initial_state {
  pieces "white" {
    A1 = "rook"
    B1 = "knight"
    C1 = "bishop"
    E1 = "queen"
    D1 = "king"
    F1 = "bishop"
    G1 = "knight"
    H1 = "rook"
    A2 = "pawn"
    B2 = "pawn"
    C2 = "pawn"
    E2 = "pawn"
    D2 = "pawn"
    F2 = "pawn"
    G2 = "pawn"
    H2 = "pawn"
  }
  pieces "black" {
    A8 = "rook"
    B8 = "knight"
    C8 = "bishop"
    E8 = "queen"
    D8 = "king"
    F8 = "bishop"
    G8 = "knight"
    H8 = "rook"
    A7 = "pawn"
    B7 = "pawn"
    C7 = "pawn"
    E7 = "pawn"
    D7 = "pawn"
    F7 = "pawn"
    G7 = "pawn"
    H7 = "pawn"
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