board {
  height = 8
  width = 8
}

// ===== PIECE TYPES SPECIFICATION =============================
// Each piece type should specify the motions it is able to perform.
//
// Motions are specified by giving a generator function name, which generates all possible destination
// squares given:
//   * the current square of the piece,
//   * the piece that is about to move.
//
// Motions can specify special actions that can alter the game state after the motion is taken via
// attribute "actions", pointing to functions that receive:
//   * the piece that moved,
//   * the starting square,
//   * the destination square,
//   * the current game state.
//  Such an action can be for example pawn promotion.
//
// Both generator and action functions are specified below the piece types definition.

piece_types {
  piece_type "king" {
    # motion {
    #   generator = "motion_castling"
    #   actions   = ["displace_rook_after_castling"]
    # }
    motion {
      generator = "motion_neighbours_straight"
    }
  }

  piece_type "queen" {
    motion {
      generator = "motion_line_diagonal"
    }
    motion {
      generator = "motion_line_straight"
    }
  }

  piece_type "rook" {
    motion {
      generator = "motion_line_straight"
    }
  }

  piece_type "knight" {
    motion {
      generator = "motion_hook"
    }
  }

  piece_type "bishop" {
    motion {
      generator = "motion_line_diagonal"
    }
  }

  piece_type "pawn" {
    motion {
      generator = "motion_forward_straight"
      # actions   = ["promote"]
    }
    # motion {
    #   generator = "motion_forward_straight_double"
    # }
    motion {
      generator = "motion_forward_diagonal"
      # actions   = ["promote"]
    }
    # motion {
    #   generator = "motion_en_passant"
    #   actions   = ["capture_en_passant"]
    # }
  }
}


// ===== MOTION GENERATOR FUNCTIONS ============================
// They receive 2 parameters:
//  * square - the current square,
//  * piece - the current piece,
// and generate all the squares that the given piece can move to from the given square.

// Generates motions to the straight neighbours (top, right, bottom, left)
// of the current square, given that they are not occupied by the player
// owning the current piece. Additionaly moving to attacked squares (ones
// which can be reached by any opponent's piece in the next turn) is blocked.
composite_function "motion_neighbours_straight" {
  params = [square, piece]
  result = {
    dposes = [[0, 1], [1, 0], [0, -1], [-1, 0]]
    dests  = [for dpos in dposes: get_square_relative(square, dpos)]
    return = [for dest in dests: dest if dest != null && !is_square_owned_by(dest, piece.color) && !is_attacked(square)]
  }
}

// Generates castling motions (both queen-side and king-side).
// Conditions:
//   * the king must have never moved in this game,
//   * the rook at the appropriate side must have never moved in this game,
//   * there must be no pieces between the king and the rook,
//   * squares on the king's path must not be attacked (including both ends).
# composite_function "motion_castling" {
#   params = [square, piece]
#   result = {
#     rooks       = [_piece for _piece in piece.owner.pieces if _piece.type_name == "rook" && !has_ever_moved(_piece)]
#     paths       = [squares_connecting_horizontal(piece.square, rook.square) for rook in rooks]
#     king_paths  = [slice(path, 0, 3) for path in paths]
#     inner_paths = [slice(path, 1, -1) for path in paths]
#     king_dests  = [path[2] for path in paths]
#     return      = [dest for i, dest in king_dests if all([!s.is_attacked for s in king_paths[i]]...) && all([s.piece == null for s in inner_paths[i]]...) && !has_ever_moved(king)]
#   }
# }

// Generates a motion one square forwards, given that the destination square
// is not occupied by any piece.
composite_function "motion_forward_straight" {
  params = [square, piece]
  result = {
    dest   = get_square_relative(square, piece.owner.forward_direction)
    return = dest != null && dest.piece == null ? [dest] : []
  }
}

// Generates a motion two square forwards, given that both the destination square and the transitional
// square are not occupied by any piece and that the piece has not moved yet before.
# composite_function "motion_forward_straight_double" {
#   params = [square, piece]
#   result = {
#     dpos   = [dcoord * 2 for dcoord in piece.owner.forward_direction]
#     dest   = get_square_relative(square, dpos)
#     middle = get_square_relative(square, piece.owner.forward_direction)
#     return = dest != null && dest.piece == null && middle.piece == null && !has_ever_moved(piece) ? [dest] : []
#   }
# }

// Generates 2 motions: one square forwards and to either side, given that the
// destination squares are occupied by the opposing player.
composite_function "motion_forward_diagonal" {
  params = [square, piece]
  result = {
    forward_y = piece.owner.forward_direction[1]
    dposes    = [[-1, forward_y], [1, forward_y]]
    dests     = [for dpos in dposes: get_square_relative(square, dpos)]
    return    = [for dest in dests: dest if dest != null && dest.piece != null && dest.piece.owner != piece.owner]
  }
}

// Generates 2 motions (en passant): one square forwards and to either side, given that the
// destination squares are free, and the last move was a "forward_straight_double"
// by an opposing pawn placed the destination file.
# composite_function "motion_en_passant" {
#   params = [square, piece]
#   result = {
#     forward   = piece.owner.forward_direction
#     forward_y = forward[1]
#     dposes    = [[-1, forward_y], [1, forward_y]]
#     dests     = [get_square_relative(square, dpos) for dpos in dposes]
#     last_move = last_or_null(game.record)
#     backward  = [-1 * dcoord for dcoord in piece.owner.forward_direction]
#     return    = [dest for dest in dests if dest != null && dest.piece == null && last_move != null && last_move.piece.type_name == "pawn" && last_move.dest == get_square_relative(dest, backward) && last_move.src == get_square_relative(dest, forward)]
#   }
# }

// Generates a maximum of 8 motions, meeting criteria:
//   * first go 2 to any side,
//   * then go 1 to any side, but the direction is perpendicular to the one of previous step.
// If the destination square is occupied by the player owing the current piece, it is discarded.
composite_function "motion_hook" {
  params = [square, piece]
  result = {
    dposes = [[2, 1], [2, -1], [-2, 1], [-2, -1], [1, 2], [-1, 2], [1, -2], [-1, -2]]
    dests  = [for dpos in dposes: get_square_relative(square, dpos)]
    return = [for dest in dests: dest if dest != null && !is_square_owned_by(dest, piece.color)]
  }
}

// Generates motions from current position (param 'square') in the given direction (param 'dpos' in
// form [dx, dy]) until end of board or a piece is encountered. If said piece belongs to the same
// player as the one in param 'piece', the last square is excluded from the generated square, else
// it is included.
composite_function "motion_line" {
  params = [square, piece, dpos]
  result = {
    next   = get_square_relative(square, dpos)
    return = next == null ? [] : next.piece == null ? list(next, motion_line(next, piece, dpos)...) : is_square_owned_by(next, piece.color) ? [] : [next]
  }
}

composite_function "motion_line_diagonal" {
  params = [square, piece]
  result = {
    dposes = [[-1, 1], [1, 1], [1, -1], [-1, -1]]
    return = concat([for dpos in dposes: motion_line(square, piece, dpos)]...)
  }
}

composite_function "motion_line_straight" {
  params = [square, piece]
  result = {
    dposes = [[0, 1], [1, 0], [0, -1], [-1, 0]]
    return = concat([for dpos in dposes: motion_line(square, piece, dpos)]...)
  }
}

initial_state {
  pieces "white" {
    A1 = "king"
  }
  pieces "black" {
    A2 = "king"
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

composite_function "decide_winner" {
  params = [game]
  result = {
    points_per_player = {for i, player in game.players: calc_player_points(player) => player...}
    best_players = points_per_player[max(keys(points_per_player)...)]
    return = length(best_players) == 1 ? best_players[0] : null
  }
}
