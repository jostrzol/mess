board {
  height = 8
  width  = 8
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
  piece_type "king" {
    representation {
      white {
        symbol = "♔"
        icon = "/piece_types/king.svg"
      }
      black {
        symbol = "♚"
        icon = "/piece_types/king.svg"
      }
    }
    motion {
      generator = "motion_castling"
      action    = "displace_rook_after_castling"
    }
    motion {
      generator = "motion_neighbours"
    }
  }

  piece_type "queen" {
    representation {
      white {
        symbol = "♕"
        icon = "/piece_types/queen.svg"
      }
      black {
        symbol = "♛"
        icon = "/piece_types/queen.svg"
      }
    }
    motion {
      generator = "motion_line_diagonal"
    }
    motion {
      generator = "motion_line_straight"
    }
  }

  piece_type "rook" {
    representation {
      white {
        symbol = "♖"
        icon = "/piece_types/rook.svg"
      }
      black {
        symbol = "♜"
        icon = "/piece_types/rook.svg"
      }
    }
    motion {
      generator = "motion_line_straight"
    }
  }

  piece_type "knight" {
    representation {
      white {
        symbol = "♘"
        icon = "/piece_types/knight.svg"
      }
      black {
        symbol = "♞"
        icon = "/piece_types/knight.svg"
      }
    }
    motion {
      generator = "motion_hook"
    }
  }

  piece_type "bishop" {
    representation {
      white {
        symbol = "♗"
        icon = "/piece_types/bishop.svg"
      }
      black {
        symbol = "♝"
        icon = "/piece_types/bishop.svg"
      }
    }
    motion {
      generator = "motion_line_diagonal"
    }
  }

  piece_type "pawn" {
    representation {
      white {
        symbol = "♙"
        icon = "/piece_types/pawn.svg"
      }
      black {
        symbol = "♟"
        icon = "/piece_types/pawn.svg"
      }
    }
    motion {
      generator       = "motion_forward_straight"
      choice_function = "promote_choose_piece_type"
      action          = "promote"
    }
    motion {
      generator = "motion_forward_straight_double"
    }
    motion {
      generator       = "motion_forward_diagonal"
      choice_function = "promote_choose_piece_type"
      action          = "promote"
    }
    motion {
      generator = "motion_en_passant"
      action    = "capture_en_passant"
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
      for dest in filternulls(dests) : dest
      if !belongs_to(piece.color, dest)
    ]
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
    rooks = [
      for _piece in owner_of(piece).pieces : _piece
      if _piece.type == "rook" && !has_ever_moved(_piece)
    ]
    paths       = [for rook in rooks : square_range(piece.square, rook.square)]
    inner_paths = [for path in paths : slice(path, 1, length(path) - 1)]
    king_dests  = [for path in paths : length(path) > 3 ? path[2] : null]
    return = has_ever_moved(piece) ? [] : [
      for i, dest in king_dests : dest
      if all([for s in inner_paths[i] : piece_at(s) == null]...)
    ]
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

// Generates a motion two square forwards, given that both the destination
// square and the transitional square are not occupied by any piece and that the
// piece has not moved yet before.
composite_function "motion_forward_straight_double" {
  params = [square, piece]
  result = {
    dpos   = [for dcoord in owner_of(piece).forward_direction : dcoord * 2]
    dest   = get_square_relative(square, dpos)
    middle = get_square_relative(square, owner_of(piece).forward_direction)
    return = dest == null ? [] : (
      piece_at(dest) != null
      || piece_at(middle) != null
      || has_ever_moved(piece)
      ? [] : [dest]
    )
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
    return = [
      for dest in filternulls(dests) : dest
      if(piece_at(dest) != null && !belongs_to(piece.color, dest)
    )]
  }
}

// Generates 2 motions (en passant): one square forwards and to either side,
// given that the destination squares are free, and the last move was a
// "forward_straight_double" by an opposing pawn placed the destination file.
composite_function "motion_en_passant" {
  params = [square, piece]
  result = {
    forward   = owner_of(piece).forward_direction
    forward_y = forward[1]
    dposes    = [[-1, forward_y], [1, forward_y]]
    dests     = [for dpos in dposes : get_square_relative(square, dpos)]
    last_move = last_or_null(game.record)
    backward  = [for dpos in owner_of(piece).forward_direction : -1 * dpos]
    return = [
      for dest in filternulls(dests) : dest
      if last_move == null ? false : (
        piece_at(dest) == null
        && last_move.piece.type == "pawn"
        && last_move.dst == get_square_relative(dest, backward)
        && last_move.src == get_square_relative(dest, forward)
      )
    ]
  }
}

// Generates a maximum of 8 motions, meeting criteria:
//   * first go 2 to any side,
//   * then go 1 to any side, but the direction is perpendicular to the one of
//     previous step.
// If the destination square is occupied by the player owing the current piece,
// it is discarded.
composite_function "motion_hook" {
  params = [square, piece]
  result = {
    dposes = [
      [2, 1], [2, -1], [-2, 1], [-2, -1],
      [1, 2], [-1, 2], [1, -2], [-1, -2]
    ]
    dests = [for dpos in dposes : get_square_relative(square, dpos)]
    return = [
      for dest in filternulls(dests) : dest
      if !belongs_to(piece.color, dest)
    ]
  }
}

// Generates motions from current position (param 'square') in the given
// direction (param 'dpos' in form [dx, dy]) until end of board or a piece is
// encountered. If said piece belongs to the same player as the one in param
// 'piece', the last square is excluded from the generated square, else it is
// included.
composite_function "motion_line" {
  params = [square, piece, dpos]
  result = {
    next = get_square_relative(square, dpos)
    return = next == null ? [] : (
      piece_at(next) == null
      ? concat([next], motion_line(next, piece, dpos))
      : belongs_to(piece.color, next) ? [] : [next]
    )
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

// Promote - choice generator
// Makes the player choose the target piece type (any except for king and pawn).
composite_function "promote_choose_piece_type" {
  params = [piece, src, dst]
  result = {
    owner     = owner_of(piece)
    forward_y = owner.forward_direction[1]
    last_y    = forward_y == 1 ? board.height - 1 : 0
    pos       = square_to_coords(dst)
    options = [
      for type in piece_types : type.name
      if !contains(["king", "pawn"], type.name)
    ]
    choice = {
      message = "Promote"
      type    = "piece_type"
      options = options
    }
    return = pos[1] == last_y ? choice : null
  }
}

// Promote - action
// Exchanges a piece for a new one, previously chosen.
composite_function "promote" {
  params = [piece, src, dst, options]
  result = {
    owner      = owner_of(piece)
    piece_type = length(options) == 0 ? null : options[0].piece_type
    _          = piece_type == null ? null : place_new_piece(piece_type.name, dst, owner.color)
  }
}

// Captures an opposing pawn after an en passant.
composite_function "capture_en_passant" {
  params = [piece, src, dst, options]
  result = {
    backward      = [for dpos in owner_of(piece).forward_direction : -1 * dpos]
    target_square = get_square_relative(dst, backward)
    target_piece  = piece_at(target_square)
    _             = capture(target_piece)
  }
}

// Displaces the rook to the appropriate square after castling.
composite_function "displace_rook_after_castling" {
  params = [piece, src, dst, options]
  result = {
    src_pos   = square_to_coords(src)
    dst_pos   = square_to_coords(dst)
    dx        = dst_pos[0] - src_pos[0]
    rook_src  = coords_to_square([dx > 0 ? board.width - 1 : 0, src_pos[1]])
    rook_dest = get_square_relative(dst, [dx > 0 ? -1 : 1, 0])
    _         = move(piece_at(rook_src), rook_dest)
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

  // If the move was performed by a king, checks if all squares on his path were
  // safe (necessary solely for castling motion).
  function "is_kings_path_save" {
    params = [move]
    result = move.name != "motion_castling" ? true : all([
      for s in square_range(move.src, move.dst)
      : !is_attacked_by(opponent_color(move.piece.color), s)
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
    pos1 = square_to_coords(end1)
    pos2 = square_to_coords(end2)
    dx   = pos2[0] - pos1[0]
    xdir = dx == 0 ? 1 : (dx) / abs(dx)
    dy   = pos2[1] - pos1[1]
    ydir = dy == 0 ? 0 : (dy) / abs(dy)
    horiz = [
      for x in range(pos1[0], pos2[0] + xdir)
      : coords_to_square([x, pos1[1]])
    ]
    vert = [
      for y in range(pos1[1] + ydir, pos2[1] + ydir)
      : coords_to_square([pos2[0], y])
    ]
    return = concat(horiz, vert)
  }
}

// Returns the given player's opponent.
function "opponent" {
  params = [player]
  result = [
    for _player in game.players : _player
    if _player.color != player.color
  ][0]
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
    A1 = "rook"
    B1 = "knight"
    C1 = "bishop"
    D1 = "queen"
    E1 = "king"
    F1 = "bishop"
    G1 = "knight"
    H1 = "rook"
    A2 = "pawn"
    B2 = "pawn"
    C2 = "pawn"
    D2 = "pawn"
    E2 = "pawn"
    F2 = "pawn"
    G2 = "pawn"
    H2 = "pawn"
  }
  black_pieces = {
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

// ===== TURN ==================================================
turn {
  choice_function = "turn_choose_move"
  action          = "turn"
}

function "turn_choose_move" {
  params = []
  result = { type = "move", message = "Choose move" }
}

composite_function "turn" {
  params = [options]
  result = {
    _ = make_move(options[0].move, slice(options, 1, length(options)))
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
    losing_player = check_mated_player(game)
    return = (
      losing_player == null
      ? [false, null]
      : [true, opponent(losing_player).color]
    )
  }
}

// Returns the check-mated player, if any - else returns null.
composite_function "check_mated_player" {
  params = [game]
  result = {
    pieces = concat([for player in game.players : player.pieces]...)
    kings  = [for piece in pieces : piece if piece.type == "king"]
    checked = [
      for king in kings : owner_of(king)
      if is_attacked_by(opponent_color(king.color), king.square)
    ]
    valid_moves = [
      for player in checked : concat(
        [for piece in player.pieces : valid_moves_for(piece)]...
      )
    ]
    mated  = [for i, player in checked : player if length(valid_moves[i]) == 0]
    return = length(mated) == 0 ? null : mated[0]
  }
}

assets = {
  piece_types = {
    "king.svg" = <<EOF
      H4sIAJhDVmUAA+VVbW+bMBD+3l/huZrUSMXYBvOShVZqWlWTsq3Suk3rNxoIYSU4MjQv/35nDIFE
      nfZ10hCJnvM99+K7s5lc71YF2qSqymUZYUYoRmk5l0leZhF+rRdWgK+vzibvbr9MH38+3KFqk6GH
      bzezj1OELdv+4Uxt+/bxFn39fo8YYbZ99xkjvKzr9di2t9st2TpEqsy+V/F6mc8rG4i2JoKRDc4Y
      I0mdYAihPUMyZRW9Yc4ppZqODWX8vDOsCmjPcre3QEfmcoXRJk+3NxLUFFHkCnjBOYJnkqSLysBG
      fN6NM5UnaJsn9TLCnAiMlmmeLetW0D4w2ut/u3Vh9z4mGVrkRRHhTawuLGtdxPtUWXNZSHV5vlgs
      RrjRW+q1SCOcbtJSJglGVa3kS9pZyfValmlZd3awy1HHsYq8TOfxOsJKvpbJ0fIvmZen6+0+mE69
      qvc6Kka1istqIdUqwqu4VvnugpLQYSEXl4ia9yAz4lGABvmUBf4ID+q1juslSiL8CXGKOFQIzTRi
      HMxCwExoLHxusN8RfOIyx8DWSPTAsBzNEm5oIDtoe9THa9BT15GTzMDA9X3khMCZIuYbKSB+iLhL
      QiMEnoecAXF2kHhjNpRc1mAeaMMGcc1wG8wEMLgJwvROuLHkMN9eH79RMdqbeANPnXs2CDrEs+Ge
      MFrF1QucS1VcnC9lkVajPxTCJAvZC6b7ERIfqut4IAZ6v4JwMZS9nq7lI+vjWndTX8LcDqbjtA/G
      g0MJ5b7eQ0Acn0EjCeehqYznQE8cIsIQMZcEjqslByaJQU0C1ummkLvwnM60deSwZkPHYY4SPUmo
      7Sk49QPfdLxJKCAiaHrOmoTAE4fZhWWdEEiCmg5DQq1u2mSmMzK2nSdOqCdO4jwdTuLhHFpwW42b
      iwHAh35ZqjzLy7Gg7xH83jq4FuuObIMI1Q8gq4Xu8SzY2UDQY4NyqEQzNKd9U+m87q5Bd3ALamy6
      vV3mdfp2eQ3huYjnL/3t1or/3Cz8Ndn/cE70bHSfN5iZif7MXp39BqubdYYZCAAA
      EOF
    "queen.svg" = <<EOF
      H4sIAJhDVmUAA7VU207cMBB95yum7gtIxJfYuW03IAEtqkRbpEKr8hY22WzabBw5KZv9+46dwC5t
      JZ6IosgzPnPOeGac+emwruGhMF2lm5QIygkUzULnVVOm5PbmgxcT6PqsybNaN0VKGk1OTw7mby6+
      nN/8uH4P3UMJ17dnVx/PgXiMfZfnjF3cXMDXb5cgqGDs/WcCZNX37YyxzWZDN5JqU7JLk7WratEx
      BDILxCCGZELQvM8JSlhmTK7p0v+E+5xzCycjZHY/jKgOYfd62Hq4Rxd6TfbPJtCqis2ZRjAHDirA
      F6UAn3leLLtx6cz7YVaaKodNlferlPg0ILAqqnLVT4blILC1XzZRsB3HvIRlVdcpecjMoee1dbYt
      jLfQtTbHb5fL5ZGtqtG/ikeEblusb9M/YvB8DrOtETJCPZfLTNDg3eSoq6b4qatmZvTvJid76bdZ
      v4I8JZ/Ax2yFD35AAxnD1TM7QVvEVCQhCE6VFHz3CLunaKyU3Qt5aO2AxriBsZYpoTJKrC+SdjMC
      ETo/ejgqyJBGaEn+zIqdhUirLQMkFo5W+Ta3hIootvwJD8ZYIV1ukbD5+EgSu3wCGf11lhjunur1
      2JGpE73Jmm6pzTol66w31XAojoG7d1yIo73auaCirqu2K2CBbcb0klARWGCzQ1QKCRh0SzsFZjst
      JuUnKU+bqsTO+GE7QEx9oYJ22M/rH5loVBA+jSeBif4FcrwxPPafk89Z+UIFvP0SKOyY6/mrlkIk
      r1gKR/5yKXb3IqEqEqA4lYGyQ4szHkR2It2YKrwNTwaOVxg4aOjGfi/wbnf5UWZuf0knB38ARa+p
      41UFAAA=
      EOF
    "rook.svg" = <<EOF
      H4sIAJhDVmUAA31TXW+bMBR976/wvJdWKrYJJm1paKV+rJrUbZWWbtrenOCAN2Ij44Tk3+8aSkKz
      dVEUzr2ce+7hEE+uN8sSraWtldEpDgnDSOq5yZTOU/w8/RCcY1Q7oTNRGi1TrA2+vjqavLv7cjv9
      8XSP6nWOnp5vHj/eIhxQ+j26pfRueoe+fntAIQkpvf+MES6cqxJKm6YhTUSMzemDFVWh5jUFIvVE
      GKIgFoYkcxmGFV4ZzOk6/cf4iDHm6bijJLNNx6qBNjObbQD3yNws8fDZQqiUbG4MkBliiMfwhVUI
      PpNMLuoOtuVsk+RWZahRmStSPCIxRoVUeeFeCq+B0db/0hcJuteY5GihyjLFa2GPg6AqxVbaYG5K
      Y0/fLxaLE5+qNb9lzzBVBflq13Pg+VrOtgSKqcRcuW0SXnrR4KC0q1Imci21ybLLTjVobSchiftG
      qbSciyqxZqWzYfOXUfp1d6mctKWCS8L7XibqQlgrtokGm313ZwQjZ4WuF8YuU9zCUjh5zE4ZiU7w
      INZKuAJlKf6EwhCFjIzg8uhx3MN4CM9e4IgNYU8YxUPYEyI2hD0h4kPYbd7le5DRbOVc/1IPXEfg
      Oup1QpAiZ94pH6Du/v+0/4q/TfytjXyfU7edd+nwAe6zfNv16ILEXgH+uygaA9ePvSp4y9n5fkPK
      y3BGWMz9DCPjc5CICW9zJTy+2FUxiS94yx1DOZj7uT8wOZxzf4yvjv4A0R68U4kEAAA=
      EOF
    "knight.svg" = <<EOF
      H4sIAJhDVmUAA3VT226bQBB971dMty+JBHvltsQkakgaVUrbSE1a9REZbNNiFq2JHefrO4uNm1QO
      gpnZuZw5M1omF0/LBtaVXdWmzYignEDVTk1Zt/OMPNx/8hMCq75oy6IxbZWR1pCL83eT91ff8vtf
      d9ewWs/h7uHy9nMOxGfsp8oZu7q/gu8/bkBQwdj1VwJk0fddythms6EbRY2dsxtbdIt6umKYyFwi
      FjEEE4KWfUmwhUNGcu0qO1IuOecunbzkLvBUV5tL85QRDhyCEF+EAnwmc5jVTZORdWFPfL9rim1l
      /alpjPU+zGazUzemNX+qMcN0HQ7c9mMONhxytg2mmK6Y1v02FWcO1P/vaB+bKq3WVWvK8myH6m/q
      sl+kgoajo6nbalp0qTWPbfnS+dvU7Wvvsu4r29So0mD0lcVqUVhbbNMWaY7eAxECvS3a1czYZUYG
      syn66oR7nKrT/UqGtXRFv4AyI19ASk9wyEFJGnpCgEqcTlB7SsMtiNDpfNAcJEpMBKkwh7DjiIGr
      z1FTBJGcagEioUHoyZCqGETkydghKk9qlFRgK0FVAEKggRFNdYg+TnkEQtJAYAHVkYtLhyy403hn
      hEYIKpXzKDeFdiqkWscOCRtF6ESF2DjnKPPxpGiiPeyABJE0x8EcLSpjT9PQNVe4DFzIzoh38QDN
      yPkityl+0LduynC3zdESyC5BBM4TkMJzY+PC4/3WB/nGFrXbcojiIzhewwfIExzy0cCh4vkNSIGc
      /xWKQ+Ew/PHAvuL51cVaFr2tn044TaIIb1bo+U7sjppGWnl+SEWMN27PY8Lm+Ge7H/f83V9TmRLg
      ewQAAA==
      EOF
    "bishop.svg" = <<EOF
      H4sIAJhDVmUAA21Ua2+bMBT93l9x531ppWL84mGWtFLTrprUbpXWblq/kUAIG8EIaB7/ftcQ8qga
      odjH9/jcc+0Lo+vNsoBVWje5KceEU0YgLWcmyctsTF6evzohgaaNyyQuTJmOSWnI9dXZ6NPtj8nz
      n6c7aFYZPL3cPHybAHFc97ecuO7t8y38/HUPnHLXvftOgCzatopcd71e07Wkps7c+zquFvmscZHo
      WiJuclGMc5q0CcEUVhnNlc34g+2CMWbppKdE003PapA2NZutgzE6M0tyXBtHlKfrG4NkBgyUhw+m
      AvyNknTe9NMOTjdRVucJrPOkXYyJoB6BRZpni3YHrAaBrf13dxLuQWOUwTwvijFZxfW541RFvE1r
      Z2YKU19+ns/nF/ZUa/MvHRimqvB8y3bgYH0dZ1sgxVTxLG+3Ef9iRZ36rUijdJWWJkn6lQOhV3U6
      2xGn3rBQ5GU6i6uoNm9lcrz41+Tl6eoyb9O6yHGI1LCWxM0irut4G5Voc1jdpyXQ1nHZzE29HJNu
      WsRtes4uGfUvyNGxZkNJ72xN39r2RGQZt3W+OcdL8zSX+hJY9xywI6kWvmR2yqkfcO6dpOrSVXG7
      gGRMHoF7lCkfJKcaHkAy6gVejyYDEppKLkAEVIcBCJ8GQr1DyJXUDzkIj/KAg1SUaw1c0ZCHIAQO
      HDgWHfjIxV5WQg5RRFoNG3lAfakOsqd4Z3Zn6Nj6674nhq77oFbMIdFVSFkYovYOoi1mnUsqlLBu
      pMAU2CI8sDEpLZXRgPugaegHFgi0rynjcgB7zR76lGnf1i0DDVimh/p4JMwDj4YqsERFNfORiDc1
      OAko99Wpy4/rGrnZETi6TB9dq+58MNUDiBCxN+BHa8hipAV9vMc4BExhnOs+roe4d4TJBy3avSTd
      e3Fi7+AIBUHidk+hnG2AEC8Mr9r2mqLK03uEqbANLNdHeLTv9fAdwapH9ut2dfYfmSKttKAFAAA=
      EOF
    "pawn.svg" = <<EOF
      H4sIAJhDVmUAA21SXW/TMBR936+4mBeQHH/FaT7WbNLaMZAQTKID8RglaRtI7cjJmnW/nuusQRMi
      suLro3PPPffay+unQwvH2vWNNTmRTBCoTWmrxuxy8rD5ECQE+qEwVdFaU+fEWHJ9dbF8s/662vy8
      v4X+uIP7h5vPn1ZAAs5/hCvO15s1fPt+B5JJzm+/ECD7YegyzsdxZGPIrNvxO1d0+6bsORK5J2IS
      RzEpWTVUBEt4ZTRn+vw/6UoI4ekEjk093tinnAgQoCNc5HU/EqUAv2VXDHvYNm2bk2Ph3gVB1xan
      2gWlba2jb7fb7XvfqbO/65lhuw57NsPMwZrIqXJyAKVYRFMoIVBMSSog0FSyOPW7BkEFS1IQTKUe
      lRjFCVUsTGAFMmZhSOWCRSAXVCYsSn2gJKoJJIkQ6ammIUs0KKYljTyGiZE/qJgJTJA0lCxKpiBF
      qY8QauSE+owrBNMzWS2meOFlVMIEFtOTOILhBKa+vDc7ucGsRfziUEUsVhQro3UUwFal9g15SSpD
      3+U5nNxP04DAT4IG2g8D/89+rqcWx2q7omyGUyYv/UUE/xzdY1tnxprn2tnLl5sIxqYa9hl2NANt
      Y+qy6DJnH031GvxlG5MdmqF2Mzod2ga3TM9YVfT7wrni5AvVM/rXCOH48Py7urr4AxpxiWIaAwAA
      EOF
  }
}
