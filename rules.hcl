// ===== BOARD SPECIFICATION ===================================
board {
  height = 8
  width  = 8
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

// ===== MOVE FUNCTIONS ========================================
// They receive 2 parameters:
//  * square - the current square,
//  * piece - the current piece,
// and generate all the squares that the given piece can move to from the given square.

// Generates moves to the straight neighbours (top, right, bottom, left)
// of the current square, given that they are not occupied by the player
// owning the current piece. Additionaly moving to checked squares (ones
// which can be reached by any opponent's piece in the next turn) is blocked also.
composite_function "move_neighbours_straight" {
  params = [square, piece]
  result = {
    dposes = [[0, 1], [1, 0], [0, -1], [-1, 0]]
    dests  = [get_square_relative(square, dpos) for dpos in dposes]
    return = [neighbour for neighbour in neighbours_straight if !is_square_owned_by(square, piece.owner) && !is_checked(square)]
  }
}

// Generates a move one square forwards, given that the destination square
// is not occupied by any piece.
composite_function "move_forward_straight" {
  params = [square, piece]
  result = {
    dest   = get_square_relative(square, piece.owner.forward_direction)
    return = dest != null && dest.piece == null ? [dest] : []
  }
}

// Generates a move two square forwards, given that both the destination square and the transitional
// square are not occupied by any piece and that the piece has not moved yet before.
composite_function "move_forward_straight_double" {
  params = [square, piece]
  result = {
    dpos   = [dcoord * 2 for dcoord in piece.owner.forward_direction]
    dest   = get_square_relative(square, dpos)
    middle = get_square_relative(square, piece.owner.forward_direction)
    return = dest != null && dest.piece == null && middle.piece == null && !has_ever_moved(piece) ? [dest] : []
  }
}

// Generates 2 moves: one square forwards and to either side, given that the
// destination squares are occupied, but not by the player owning the current piece.
composite_function "move_forward_diagonal" {
  params = [square, piece]
  result = {
    forward_y = piece.owner.forward_direction[1]
    dposes    = [[-1, forward_y], [1, forward_y]]
    dests     = [get_square_relative(square, dpos) for dpos in dposes]
    return    = [dest for dest in dests if dest != null && dest.piece != null && dest.piece.owner != piece.owner]
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
    dests  = [get_square_relative(square, dpos) for dpos in dposes]
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
    next   = get_square_relative(square, dpos)
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

// ===== PIECE TYPES SPECIFICATION =============================
// Each piece type should specify the moves it should be able to perform.
//
// Moves are specified by giving a generator function, which generates all possible destination
// squares given:
//   * the current square of the piece,
//   * the piece that is about to move.
//
// Moves can specify special actions that can alter the game state after the move is taken via
// attribute "action", pointing to a function that receives:
//   * the piece that moved,
//   * the starting square,
//   * the destination square,
//   * the current game state.
//  Such an action can be for example pawn promotion.

piece_types {
  piece_type "king" {
    // TODO: castling
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
    // TODO: promotion
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
// Returns the other player than the one given.
function "other_player" {
  params = [this_player]
  result = [player for player in game.players if player != this_player]
}

// Returns the check_mated king, if any - else returns null.
function "check_mated_king" {
  params = [game]
  result = {
    kings   = [piece for piece in game.players.*.pieces if piece.type_name == "king"]
    checked = [king for king in kings if is_checked(king.square)]
    mated   = [king for king in checked if length(valid_moves(king)) == 0]
    return  = length(mated) == 0 ? null : mated[0]
  }
}

// This function is called at the start of every turn.
// Returns a tuple in form [is_finished, winner]. If is_finished == true and winner == null
// then draw is concluded. Stalemate is "hardcoded" into the game - the rules don't have
// to specify it explicitly.
function "pick_winner" {
  params = [game]
  result = {
    losers_king = check_mated_king(game)
    return      = losers_king == null ? [false, null] : [true, other_player(losers_king.owner)]
  }
}