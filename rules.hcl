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
    motion {
      generator = "motion_forward_straight_double"
    }
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
// owning the current piece.
composite_function "motion_neighbours_straight" {
  params = [square, piece]
  result = {
    dposes = [[0, 1], [1, 0], [0, -1], [-1, 0]]
    dests  = [for dpos in dposes: get_square_relative(square, dpos)]
    return = [for dest in dests: dest if dest == null ? false : !belongs_to(piece.color, dest)]
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
#     rooks       = [_piece for _piece in owner_of(piece).pieces if _piece.type == "rook" && !has_ever_moved(_piece)]
#     paths       = [square_range(piece.square, rook.square) for rook in rooks]
#     inner_paths = [slice(path, 1, -1) for path in paths]
#     king_dests  = [path[2] for path in paths]
#     return      = [dest for i, dest in king_dests if all([s.piece == null for s in inner_paths[i]]...) && !has_ever_moved(king)]
#   }
# }

// Generates a motion one square forwards, given that the destination square
// is not occupied by any piece.
composite_function "motion_forward_straight" {
  params = [square, piece]
  result = {
    dest   = get_square_relative(square, owner_of(piece).forward_direction)
    return = dest == null ? [] : piece_at(dest) != null ? [] : [dest]
  }
}

// Generates a motion two square forwards, given that both the destination square and the transitional
// square are not occupied by any piece and that the piece has not moved yet before.
composite_function "motion_forward_straight_double" {
  params = [square, piece]
  result = {
    dpos   = [for dcoord in owner_of(piece).forward_direction: dcoord * 2]
    dest   = get_square_relative(square, dpos)
    middle = get_square_relative(square, owner_of(piece).forward_direction)
    return = dest == null ? [] : piece_at(dest) != null || piece_at(middle) != null || has_ever_moved(piece) ? [] : [dest]
  }
}

// Generates 2 motions: one square forwards and to either side, given that the
// destination squares are occupied by the opposing player.
composite_function "motion_forward_diagonal" {
  params = [square, piece]
  result = {
    forward_y = owner_of(piece).forward_direction[1]
    dposes    = [[-1, forward_y], [1, forward_y]]
    dests     = [for dpos in dposes: get_square_relative(square, dpos)]
    return    = [for dest in dests: dest if dest == null ? false : piece_at(dest) != null && !belongs_to(piece.color, dest)]
  }
}

// Generates 2 motions (en passant): one square forwards and to either side, given that the
// destination squares are free, and the last move was a "forward_straight_double"
// by an opposing pawn placed the destination file.
# composite_function "motion_en_passant" {
#   params = [square, piece]
#   result = {
#     forward   = owner_of(piece).forward_direction
#     forward_y = forward[1]
#     dposes    = [[-1, forward_y], [1, forward_y]]
#     dests     = [get_square_relative(square, dpos) for dpos in dposes]
#     last_move = last_or_null(game.record)
#     backward  = [-1 * dcoord for dcoord in owner_of(piece).forward_direction]
#     return    = [dest for dest in dests if dest == null ? false : piece_at(dest) == null && last_move != null && last_move.piece.type == "pawn" && last_move.dest == get_square_relative(dest, backward) && last_move.src == get_square_relative(dest, forward)]
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
    return = [for dest in dests: dest if dest == null ? false : !belongs_to(piece.color, dest)]
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
    return = next == null ? [] : piece_at(next) == null ? concat([next], motion_line(next, piece, dpos)) : belongs_to(piece.color, next) ? [] : [next]
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

// ===== GAME STATE VALIDATORS ===================================
// Validators are called just after a move is taken. If any validator returns false, then the move
// is reversed - it cannot be completed.
// Validators receive:
//   * the current game state,
//   * the last move,
// and return true if the state is valid or false otherwise.


state_validators {
  // Checks if the current player's king is not standing on an attacke square
  function "is_king_safe" {
    params = [game, move]
    result = all([for piece in move.player.pieces: !is_attacked_by(opponent_color(piece.color), piece.square) if piece.type == "king"]...)
  }

  // If the last move was performed by a king, checks if all squares on his path were safe
  // function "is_kings_path_save" {
  //   params = [game, move]
  //   result = move.piece.type != "king" ? true : all([for s in square_range(move.src, move.dst): !is_attacked(s)]...)
  // }
}

// ===== HELPER FUNCTIONS ======================================
// Checks if square is occupied by a piece of a given color
composite_function "belongs_to" {
  params = [color, square]
  result = {
    piece = piece_at(square)
    return = piece == null ? false : piece.color == color
  }
}

// Checks if the given piece has ever moved in the current game.
function "has_ever_moved" {
  params = [piece]
  result = length([for move in game.record: move if move.piece == piece]) != 0
}

// Returns the last element in the given collection or null if empty.
function "last_or_null" {
  params = [collection]
  result = length(collection) == 0 ? null : collection[length(collection) - 1]
}

// Returns all the squares connecting two given end-squares, forming an L-shape
// (including the end-squares)
composite_function "square_range" {
  params = [end1, end2]
  result = {
    pos1 = square_to_coords(end1)
    pos2 = square_to_coords(end2)
    horiz = [for x in range(pos1[0], pos2[0] + 1): coords_to_square([x, pos1[1]])]
    vert = [for y in range(pos1[1], pos2[1] + 1): coords_to_square([pos1[0], y])]
    return = concat(horiz, vert)
  }
}

initial_state {
  pieces "white" {
    A3 = "king"
    B3 = "bishop"
    C3 = "rook"
    D3 = "knight"
    E3 = "queen"
    C1 = "rook"
    E1 = "pawn"
    F1 = "pawn"
  }
  pieces "black" {
    A1 = "king"
    B1 = "queen"
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

// ===== GAME RESOLVING FUNCTIONS ==============================
// Namely the function "pick_winner" and its helpers

// This function is called at the start of every turn.
// Returns a tuple in form [is_finished, winner_color]. If is_finished == true and
// winner_color == null then draw is concluded. Stalemate is "hardcoded" into the game - the rules
// don't have to specify it explicitly.
composite_function "pick_winner" {
  params = [game]
  result = {
    losing_player = check_mated_player(game)
    return        = losing_player == null ? [false, null] : [true, opponent(losing_player).color]
  }
}

// Returns the check-mated player, if any - else returns null.
composite_function "check_mated_player" {
  params = [game]
  result = {
    pieces  = concat([for player in game.players: player.pieces]...)
    kings   = [for piece in pieces: piece if piece.type == "king"]
    checked = [for king in kings: owner_of(king) if is_attacked_by(opponent_color(king.color), king.square)]
    valid_moves = [for player in checked: concat([for piece in player.pieces: valid_moves_for(piece)]...)]
    mated   = [for i, player in checked: player if length(valid_moves[i]) == 0]
    return  = length(mated) == 0 ? null : mated[0]
  }
}

// Returns the given player's opponent.
function "opponent" {
  params = [player]
  result = [for _player in game.players: _player if _player.color != player.color][0]
}

// Returns the color belonging to the opponent of the player having the given color.
function "opponent_color" {
  params = [color]
  result = [for _player in game.players: _player.color if _player.color != color][0]
}
