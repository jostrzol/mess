board {
  height = 8
  width  = 8
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
    motion {
      generator = "motion_castling"
      actions   = ["displace_rook_after_castling"]
    }
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
    motion {
      generator = "motion_en_passant"
      actions   = ["capture_en_passant"]
    }
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
    dests  = [for dpos in dposes : get_square_relative(square, dpos)]
    return = [for dest in dests : dest if dest == null ? false : !belongs_to(piece.color, dest)]
  }
}

// Generates castling motions (both queen-side and king-side).
// Conditions:
//   * the king must have never moved in this game,
//   * the rook at the appropriate side must have never moved in this game,
//   * there must be no pieces between the king and the rook,
//   * squares on the king's path must not be attacked (including both ends).
composite_function "motion_castling" {
  params = [square, piece]
  result = {
    rooks       = [for _piece in owner_of(piece).pieces : _piece if _piece.type == "rook" && !has_ever_moved(_piece)]
    paths       = [for rook in rooks : square_range(piece.square, rook.square)]
    inner_paths = [for path in paths : slice(path, 1, length(path) - 1)]
    king_dests  = [for path in paths : length(path) > 3 ? path[2] : null]
    return      = has_ever_moved(piece) ? [] : [for i, dest in king_dests : dest if all([for s in inner_paths[i] : piece_at(s) == null]...)]
  }
}

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
    dpos   = [for dcoord in owner_of(piece).forward_direction : dcoord * 2]
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
    dests     = [for dpos in dposes : get_square_relative(square, dpos)]
    return    = [for dest in dests : dest if dest == null ? false : piece_at(dest) != null && !belongs_to(piece.color, dest)]
  }
}

// Generates 2 motions (en passant): one square forwards and to either side, given that the
// destination squares are free, and the last move was a "forward_straight_double"
// by an opposing pawn placed the destination file.
composite_function "motion_en_passant" {
  params = [square, piece]
  result = {
    forward   = owner_of(piece).forward_direction
    forward_y = forward[1]
    dposes    = [[-1, forward_y], [1, forward_y]]
    dests     = [for dpos in dposes : get_square_relative(square, dpos)]
    last_move = last_or_null(game.record)
    backward  = [for dpos in owner_of(piece).forward_direction : -1 * dpos]
    return    = [for dest in dests : dest if dest == null || last_move == null ? false : piece_at(dest) == null && last_move.piece.type == "pawn" && last_move.dst == get_square_relative(dest, backward) && last_move.src == get_square_relative(dest, forward)]
  }
}

// Generates a maximum of 8 motions, meeting criteria:
//   * first go 2 to any side,
//   * then go 1 to any side, but the direction is perpendicular to the one of previous step.
// If the destination square is occupied by the player owing the current piece, it is discarded.
composite_function "motion_hook" {
  params = [square, piece]
  result = {
    dposes = [[2, 1], [2, -1], [-2, 1], [-2, -1], [1, 2], [-1, 2], [1, -2], [-1, -2]]
    dests  = [for dpos in dposes : get_square_relative(square, dpos)]
    return = [for dest in dests : dest if dest == null ? false : !belongs_to(piece.color, dest)]
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
    return = concat([for dpos in dposes : motion_line(square, piece, dpos)]...)
  }
}

composite_function "motion_line_straight" {
  params = [square, piece]
  result = {
    dposes = [[0, 1], [1, 0], [0, -1], [-1, 0]]
    return = concat([for dpos in dposes : motion_line(square, piece, dpos)]...)
  }
}

// ===== ACTION FUNCTIONS ========================================
// Actions executed after piece movement.
// They receive 3 parameters:
// * the piece that moved,
// * the source square,
// * the destination square.

// Exchanges a piece for a new one of any type except pawn and king. Works only if moved to the last
// rank.
# composite_function "promote" {
#   params = [piece, src, dst]
#   result = {
#     valid_piece_types = [type for type in piece_types if !contains(["king", "pawn"], type.name)]
#     forward_y         = piece.owner.forward_direction[1]
#     last_rank         = forward_y == 1 ? board.height : 1
#     _                 = dest.rank == last_rank ? exchange_piece(piece, valid_piece_types) : null
#     return            = null
#   }
# }

// Captures an opposing pawn after an en passant.
composite_function "capture_en_passant" {
  params = [piece, src, dst]
  result = {
    backward        = [for dpos in owner_of(piece).forward_direction : -1 * dpos]
    attacked_square = get_square_relative(dst, backward)
    attacked_piece  = piece_at(attacked_square)
    _               = capture(attacked_piece)
    return          = null
  }
}

// Displaces the rook to the appropriate square after castling.
composite_function "displace_rook_after_castling" {
  params = [piece, src, dst]
  result = {
    src_pos   = square_to_coords(src)
    dst_pos   = square_to_coords(dst)
    dx        = dst_pos[0] - src_pos[0]
    rook_src  = coords_to_square([dx > 0 ? board.width - 1 : 0, src_pos[1]])
    rook_dest = get_square_relative(dst, [dx > 0 ? -1 : 1, 0])
    _         = move(piece_at(rook_src), rook_dest)
    return    = null
  }
}

// ===== GAME STATE VALIDATORS ===================================
// Validators are called just after a move is taken. If any validator returns false, then the move
// is reversed - it cannot be completed.
// Validators receive 1 parameter - the last move and return true if the state is valid or false otherwise.

state_validators {
  // Checks if the current player's king is not standing on an attacked square
  function "is_king_safe" {
    params = [move]
    result = all([for piece in move.player.pieces : !is_attacked_by(opponent_color(piece.color), piece.square) if piece.type == "king"]...)
  }

  // If the move was performed by a king, checks if all squares on his path were safe
  function "is_kings_path_save" {
    params = [move]
    result = move.piece.type != "king" ? true : all([for s in square_range(move.src, move.dst) : !is_attacked_by(opponent_color(move.piece.color), s)]...)
  }
}

// ===== HELPER FUNCTIONS ======================================
// Checks if square is occupied by a piece of a given color
composite_function "belongs_to" {
  params = [color, square]
  result = {
    piece  = piece_at(square)
    return = piece == null ? false : piece.color == color
  }
}

// Checks if the given piece has ever moved in the current game.
function "has_ever_moved" {
  params = [piece]
  result = length([for move in game.record : move if move.piece == piece]) != 0
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
    pos1   = square_to_coords(end1)
    pos2   = square_to_coords(end2)
    xdir   = pos2[0] - pos1[0] == 0 ? 1 : (pos2[0] - pos1[0]) / abs(pos2[0] - pos1[0])
    ydir   = pos2[1] - pos1[1] == 0 ? 0 : (pos2[1] - pos1[1]) / abs(pos2[1] - pos1[1])
    horiz  = [for x in range(pos1[0], pos2[0] + xdir) : coords_to_square([x, pos1[1]])]
    vert   = [for y in range(pos1[1] + ydir, pos2[1] + ydir) : coords_to_square([pos2[0], y])]
    return = concat(horiz, vert)
  }
}

// Returns the given player's opponent.
function "opponent" {
  params = [player]
  result = [for _player in game.players : _player if _player.color != player.color][0]
}

// Returns the color belonging to the opponent of the player having the given color.
function "opponent_color" {
  params = [color]
  result = [for _player in game.players : _player.color if _player.color != color][0]
}

// ===== INITIAL STATE =========================================
initial_state {
  pieces "white" {
    E1 = "king"
    A1 = "rook"
    H1 = "rook"
    G1 = "knight"
    B5 = "pawn"
  }
  pieces "black" {
    A8 = "rook"
    B8 = "knight"
    C8 = "bishop"
    D8 = "queen"
    E8 = "king"
    F8 = "bishop"
    G8 = "knight"
    H8 = "rook"
    A7 = "pawn"
    B7 = "pawn"
    C7 = "pawn"
    D7 = "pawn"
    E7 = "pawn"
    F7 = "pawn"
    G7 = "pawn"
    H7 = "pawn"
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
    pieces      = concat([for player in game.players : player.pieces]...)
    kings       = [for piece in pieces : piece if piece.type == "king"]
    checked     = [for king in kings : owner_of(king) if is_attacked_by(opponent_color(king.color), king.square)]
    valid_moves = [for player in checked : concat([for piece in player.pieces : valid_moves_for(piece)]...)]
    mated       = [for i, player in checked : player if length(valid_moves[i]) == 0]
    return      = length(mated) == 0 ? null : mated[0]
  }
}
