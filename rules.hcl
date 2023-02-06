// ===== BOARD SPECIFICATION ===================================
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
      actions   = ["promote"]
    }
    motion {
      generator = "motion_forward_straight_double"
    }
    motion {
      generator = "motion_forward_diagonal"
      actions   = ["promote"]
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
// owning the current piece. Additionaly moving to attacked squares (ones
// which can be reached by any opponent's piece in the next turn) is blocked.
composite_function "motion_neighbours_straight" {
  params = [square, piece]
  result = {
    dposes = [[0, 1], [1, 0], [0, -1], [-1, 0]]
    dests  = [get_square_relative(square, dpos) for dpos in dposes]
    return = [neighbour for neighbour in neighbours_straight if !is_square_owned_by(square, piece.owner) && !is_attacked(square)]
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
    rooks       = [_piece for _piece in piece.owner.pieces if _piece.type_name == "rook" && !has_ever_moved(_piece)]
    paths       = [squares_connecting_horizontal(piece.square, rook.square) for rook in rooks]
    king_paths  = [slice(path, 0, 3) for path in paths]
    inner_paths = [slice(path, 1, -1) for path in paths]
    king_dests  = [path[2] for path in paths]
    return      = [dest for i, dest in king_dests if all([!s.is_attacked for s in king_paths[i]]...) && all([s.piece == null for s in inner_paths[i]]...) && !has_ever_moved(king)]
  }
}

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
composite_function "motion_forward_straight_double" {
  params = [square, piece]
  result = {
    dpos   = [dcoord * 2 for dcoord in piece.owner.forward_direction]
    dest   = get_square_relative(square, dpos)
    middle = get_square_relative(square, piece.owner.forward_direction)
    return = dest != null && dest.piece == null && middle.piece == null && !has_ever_moved(piece) ? [dest] : []
  }
}

// Generates 2 motions: one square forwards and to either side, given that the
// destination squares are occupied by the opposing player.
composite_function "motion_forward_diagonal" {
  params = [square, piece]
  result = {
    forward_y = piece.owner.forward_direction[1]
    dposes    = [[-1, forward_y], [1, forward_y]]
    dests     = [get_square_relative(square, dpos) for dpos in dposes]
    return    = [dest for dest in dests if dest != null && dest.piece != null && dest.piece.owner != piece.owner]
  }
}

// Generates 2 motions (en passant): one square forwards and to either side, given that the
// destination squares are free, and the last move was a "forward_straight_double"
// by an opposing pawn placed the destination file.
composite_function "motion_en_passant" {
  params = [square, piece]
  result = {
    forward   = piece.owner.forward_direction
    forward_y = forward[1]
    dposes    = [[-1, forward_y], [1, forward_y]]
    dests     = [get_square_relative(square, dpos) for dpos in dposes]
    last_move = last_or_null(game.record)
    backward  = [-1 * dcoord for dcoord in piece.owner.forward_direction]
    return    = [dest for dest in dests if dest != null && dest.piece == null && last_move != null && last_move.piece.type_name == "pawn" && last_move.dest == get_square_relative(dest, backward) && last_move.src == get_square_relative(dest, forward)]
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
    dests  = [get_square_relative(square, dpos) for dpos in dposes]
    return = [dest for dest in dests if dest != null && !is_square_owned_by(piece.owner)]
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
    return = next == null ? [] : next.piece == null ? [next, motion_line(next, piece, dpos)...] : is_square_owned_by(next, piece.owner) ? [] : [next]
  }
}

composite_function "motion_line_diagonal" {
  params = [square, piece]
  result = {
    dposes = [[-1, 1], [1, 1], [1, -1], [-1, -1]]
    return = [motion_line(square, piece, dpos)... for dpos in dposes]
  }
}

composite_function "motion_line_straight" {
  params = [square, piece]
  result = {
    dposes = [[0, 1], [1, 0], [0, -1], [-1, 0]]
    return = [motion_line(square, piece, dpos)... for dpos in dposes]
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
  function "is_king_safe" {
    params = [game, move]
    result = all([!is_attacked(piece.square) for piece in move.player.pieces if piece.type_name == "king"]...)
  }
}

// ===== ACTION FUNCTIONS ========================================
// Actions executed after piece movement. They receive 4 parameters:
//  * piece - the piece that just moved,
//  * src - the source square,
//  * dest - the destination square,
//  * game - the current state,
// and can alter the game's state.

// Exchanges a piece for a new one of any type except pawn and king. Works only if moved to the last
// rank.
composite_function "promote" {
  params = [piece, src, dest, game]
  result = {
    valid_piece_types = [type for type in piece_types if !contains(["king", "pawn"], type.name)]
    forward_y         = piece.owner.forward_direction[1]
    last_rank         = forward_y == 1 ? board.height : 1
    _                 = dest.rank == last_rank ? exchange_piece(piece, valid_piece_types) : null
    return            = null
  }
}

// Captures an opposing pawn after an en passant.
composite_function "capture_en_passant" {
  params = [piece, src, dest, game]
  result = {
    backward         = [-1 * dcoord for dcoord in piece.owner.forward_direction]
    piece_to_capture = get_square_relative(dest, backward).piece
    _                = capture(piece_to_capture)
    return           = null
  }
}

// Displaces the rook to the appropriate square after castling.
composite_function "displace_rook_after_castling" {
  params = [piece, src, dest, game]
  result = {
    dx         = dest.position[0] - src.position[0]
    rook_src_x = dx > 0 ? board.width
    rook_src   = get_square_absolute([rook_src_x, src.position[1]])
    rook_dest  = get_square_relative(dest, [dx > 0 ? -1 : 1, 0])
    _          = move(rook_square.piece, rook_dest)
    return     = null
  }
}

// ===== HELPER FUNCTIONS ======================================
// Checks if square is occupied by a piece of a given player.
function "is_square_owned_by" {
  params = [square, player]
  result = square.piece != null && square.piece.owner == player
}

// Checks if the given piece has ever moved in the current game.
function "has_ever_moved" {
  params = [piece]
  result = length([move for move in game.record if move.piece == piece]) != 0
}

// Returns the last element in the given collection or null if empty.
function "last_or_null" {
  params = [collection]
  result = length(collection) == 0 ? null : collection[length(collection) - 1]
}

// Returns all the squares connecting two given end-squares (including the end-squares)
function "squares_connecting_horizontal" {
  params = [end1, end2]
  result = [get_square_absolute([x, end1.position[1]]) for x in range(end1.position[0], end2.position[0] + 1)]
}

// ===== INITIAL GAME STATE SPECIFICATION ======================
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

// ===== GAME RESOLVING FUNCTIONS ==============================
// Namely the function "pick_winner" and its helpers

// This function is called at the start of every turn.
// Returns a tuple in form [is_finished, winner]. If is_finished == true and winner == null
// then draw is concluded. Stalemate is "hardcoded" into the game - the rules don't have
// to specify it explicitly.
composite_function "pick_winner" {
  params = [game]
  result = {
    losing_king = check_mated_king(game)
    return      = losing_king == null ? [false, null] : [true, other_player(losing_king.owner)]
  }
}

// Returns the check_mated king, if any - else returns null.
composite_function "check_mated_king" {
  params = [game]
  result = {
    kings   = [piece for piece in game.players.*.pieces if piece.type_name == "king"]
    checked = [king for king in kings if is_attacked(king.square)]
    mated   = [king for king in attacked if length(valid_moves(king)) == 0]
    return  = length(mated) == 0 ? null : mated[0]
  }
}

// Returns the other player than the one given.
function "other_player" {
  params = [this_player]
  result = [player for player in game.players if player != this_player]
}
