board {
  height = 4
  width  = 3
}

// ===== PIECE TYPES SPECIFICATION =============================
// Each piece type should specify the motions it is able to perform.
//
// Motions are specified by giving a generator function name, which generates
// all possible destination squares given:
//   * the current square of the piece,
//   * the piece that is about to move.
//
// Motions can specify special actions that can alter the game state after the
// motion is taken via attribute "actions", pointing to functions that receive:
//   * the piece that moved,
//   * the starting square,
//   * the destination square,
//   * the current game state.
//  Such an action can be for example pawn promotion.
//
// Both generator and action functions are specified below the piece types
// definition.

piece_types {
  piece_type "lion" {
    motion {
      generator = "motion_neighbours"
    }
  }

  piece_type "hen" {
    motion {
      generator = "motion_hen"
    }
  }

  piece_type "giraffe" {
    motion {
      generator = "motion_neighbours_straight"
    }
  }

  piece_type "elephant" {
    motion {
      generator = "motion_neighbours_diagonal"
    }
  }

  piece_type "chick" {
    motion {
      generator = "motion_forward_straight"
      actions   = ["promote"]
    }
  }
}


// ===== MOTION GENERATOR FUNCTIONS ============================
// They receive 2 parameters:
//  * square - the current square,
//  * piece - the current piece,
// and generate all the squares that the given piece can move to from the given
// square.

// Generates motions to all the 8 neighbours of the current square,
// given that they are not occupied by the player owning the current piece.
composite_function "motion_neighbours" {
  params = [square, piece]
  result = {
    dposes = [
      [0, 1], [1, 0], [0, -1], [-1, 0],
      [1, 1], [1, -1], [-1, 1], [-1, -1]
    ]
    dests = [for dpos in dposes : get_square_relative(square, dpos)]
    return = [
      for dest in dests : dest
      if dest == null ? false : !belongs_to(piece.color, dest)
    ]
  }
}


// Generates motions to all the 4 side-neighbours of the current square + 2 forward
// diagonal-neighbours, given that they are not occupied by the player owning the current piece.
composite_function "motion_hen" {
  params = [square, piece]
  result = {
    forward_y = owner_of(piece).forward_direction[1]
    dposes = [
      [-1, forward_y], [1, forward_y],
      [0, 1], [1, 0], [0, -1], [-1, 0]
    ]
    dests = [for dpos in dposes : get_square_relative(square, dpos)]
    return = [
      for dest in dests : dest
      if dest == null ? false : !belongs_to(piece.color, dest)
    ]
  }
}

// Generates motions to all the 4 side-neighbours of the current square,
// given that they are not occupied by the player owning the current piece.
composite_function "motion_neighbours_straight" {
  params = [square, piece]
  result = {
    dposes = [[0, 1], [1, 0], [0, -1], [-1, 0]]
    dests  = [for dpos in dposes : get_square_relative(square, dpos)]
    return = [
      for dest in dests : dest
      if dest == null ? false : !belongs_to(piece.color, dest)
    ]
  }
}

// Generates motions to all the 4 corner-neighbours of the current square,
// given that they are not occupied by the player owning the current piece.
composite_function "motion_neighbours_diagonal" {
  params = [square, piece]
  result = {
    dposes = [[1, 1], [1, -1], [-1, 1], [-1, -1]]
    dests  = [for dpos in dposes : get_square_relative(square, dpos)]
    return = [
      for dest in dests : dest
      if dest == null ? false : !belongs_to(piece.color, dest)
    ]
  }
}

// Generates a motion one square forwards, given that the destination square
// is not occupied by own piece.
composite_function "motion_forward_straight" {
  params = [square, piece]
  result = {
    dest   = get_square_relative(square, owner_of(piece).forward_direction)
    return = dest == null ? [] : belongs_to(piece.color, dest) ? [] : [dest]
  }
}

// ===== ACTION FUNCTIONS ========================================
// Actions executed after piece movement.
// They receive 3 parameters:
// * the piece that moved,
// * the source square,
// * the destination square.

// Exchanges a piece for a hen.
composite_function "promote" {
  params = [piece, src, dst]
  result = {
    owner     = owner_of(piece)
    forward_y = owner.forward_direction[1]
    last_y    = forward_y == 1 ? board.height - 1 : 0
    pos       = square_to_coords(dst)
    return    = cond_call(pos[1] == last_y, "place_new_piece", "hen", dst, owner.color)
  }
}

// ===== GAME STATE VALIDATORS ===================================
// Validators are called just after a move is taken. If any validator returns
// false, then the move is reversed - it cannot be completed.
//
// Validators receive 1 parameter - the last move and return true if the state
// is valid or false otherwise.

state_validators {
  // Checks if the current player's king is not standing on an attacked square
  function "is_king_safe" {
    params = [move]
    result = all([
      for piece in move.player.pieces
      : !is_attacked_by(opponent_color(piece.color), piece.square)
      if piece.type == "king"
    ]...)
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

// Returns the color belonging to the opponent of the player having the given
// color.
function "opponent_color" {
  params = [color]
  result = [
    for _player in game.players : _player.color
    if _player.color != color
  ][0]
}

// ===== INITIAL STATE =========================================
initial_state {
  white_pieces = {
    A1 = "elephant"
    B1 = "lion"
    C1 = "giraffe"
    B2 = "chick"
  }
  black_pieces = {
    A4 = "giraffe"
    B4 = "lion"
    C4 = "elephant"
    B3 = "chick"
  }
}

// ===== TURN ==================================================
composite_function "turn" {
  params = []
  result = {
    func   = length(game.current_player.captures) == 0 ? "player_move" : "choose_turn_action"
    _      = call(func)
  }
}

composite_function "choose_turn_action" {
  params = []
  result = {
    actions      = ["Move piece", "Place captured piece"]
    action_funcs = ["player_move", "place_capture"]
    choice       = choose(actions)
    _            = call(action_funcs[choice])
  }
}

composite_function "place_capture" {
  params = []
  result = {
    captures = [for type, count in game.current_player.captures : [type, count]]
    types    = [for pair in captures : format("%s (count: %v)", pair[0], pair[1])]
    type_idx = choose(types)
    type     = captures[type_idx][0]
    coords_list = concat(
      [
        for x in range(board.width)
        : [for y in range(board.height) : [x, y]]
      ]...
    )
    squares       = [for coords in coords_list : coords_to_square(coords)]
    empty_squares = [for square in squares : square if piece_at(square) == null]
    square_idx    = choose(empty_squares)
    square        = empty_squares[square_idx]
    _             = convert_and_release(game.current_player, type, square)
  }
}

// ===== GAME RESOLVING FUNCTIONS ==============================
// Namely the function "pick_winner" and its helpers

// This function is called at the start of every turn.
// Returns a tuple in form [is_finished, winner_color]. If is_finished == true
// and winner_color == null then draw is concluded. Stalemate is "hardcoded"
// into the game - the rules don't have to specify it explicitly.
composite_function "pick_winner" {
  params = [game]
  result = {
    winners = concat(
      [for player in game.players : opponent_color(player.color) if has_lost(player)],
      [for player in game.players : player.color if has_won(player)]
    )
    return = (
      length(winners) > 0
      ? [true, winners[0]]
      : [false, null]
    )
  }
}

// Checks if the given player has lost.
composite_function "has_lost" {
  params = [player]
  result = {
    lions  = [for piece in player.pieces : piece if piece.type == "lion"]
    return = length(lions) == 0
  }
}

// Checks if the given player has won.
composite_function "has_won" {
  params = [player]
  result = {
    lions  = [for piece in player.pieces : piece if piece.type == "lion"]
    return = any([for lion in lions : is_in_final_rank(player.forward_direction, lion)]...)
  }
}

// Checks if the given piece is in the final rank.
composite_function "is_in_final_rank" {
  params = [forward_direction, piece]
  result = {
    forward_y = forward_direction[1]
    last_y    = forward_y == 1 ? board.height - 1 : 0
    pos       = square_to_coords(piece.square)
    return    = pos[1] == last_y
  }
}