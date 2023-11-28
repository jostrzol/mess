// ===== BOARD ================================================================
// Board size definition.
board {
  height = 8
  width  = 8
}

// ===== PIECE TYPES SPECIFICATION ============================================
// Each piece type should specify the motions it is able to perform.
//
// Motions are specified by giving a generator function name, which generates
// all possible destination squares given:
//   * the current square of the piece,
//   * the piece that is about to move.
//
// Motions can specify special action, such as pawn promotion in chess, that
// can alter the game state after the motion is taken. It can be defined via
// the attribute named "action", which points to a function receiving:
//   * the piece that moved,
//   * the starting square,
//   * the destination square,
//   * (optionally) user options.
// Actions are not expected to produce any result, but can use builtin
// functions to modify the game state.
//
// The last argument to an action function contains the user's decisions
// regarding the current action. Choice tree, from which such decisions
// can be made, is specified via another attribute of motion configuration --
// "choice_function". This function receives arguments:
//   * the piece that moved,
//   * the starting square,
//   * the destination square,
// and is expected to return a choice tree object.
//
// Generator, action and choice functions are implemented below the piece types
// definition.
//
// Piece appearance can be also configured via the "presentation" block inside
// the piece type's definition.

piece_types {
  piece_type "piece" {
    presentation {
      white {
        symbol = "○"
        icon = "/piece_types/disk.svg"
      }
      black {
        symbol = "●"
        icon = "/piece_types/disk.svg"
      }
    }
    motion {
      generator = "motion_neighbours"
    }
    motion {
      generator = "motion_jump"
    }
  }
}

// ===== MOTION GENERATOR FUNCTIONS ===========================================
// They receive 2 parameters:
//  * square - the current square,
//  * piece - the current piece,
// and generate list of squares that the given piece can move to from the given
// square.

// Generates motions to all the 8 neighbours of the current square,
// given that they are not occupied.
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
      if dest == null ? false : !is_occupied(dest)
    ]
  }
}

// Generates motions from current position by jumping over a piece in any
// direction iteratively, until end of board or a blocking piece is encountered.
composite_function "motion_jump" {
  params = [square, piece]
  result = {
    return = motion_jump_step([], [square])
  }
}

// Iteration step of motion jump.
composite_function "motion_jump_step" {
  params = [result, to_process]
  result = {
    curr_square = to_process[0]
    dests       = motion_jump_step_impl(curr_square)
    deduped_dests = [
      for square in dests : square
      if !contains(result, square)
    ]
    new_result = concat(result, deduped_dests)
    new_to_process = (
      length(to_process) == 1
      ? deduped_dests
      : concat(
        slice(to_process, 1, length(to_process) - 1),
        deduped_dests
      )
    )
    next_step_result = cond_call(
      length(new_to_process) > 0,
      "motion_jump_step",
      new_result,
      new_to_process
    )
    return = (
      next_step_result == null
      ? new_result
      : next_step_result
    )
  }
}

// Implementation of motion jump step.
composite_function "motion_jump_step_impl" {
  params = [square]
  result = {
    mid_dposes = [
      [0, 1], [1, 0], [0, -1], [-1, 0],
      [1, 1], [1, -1], [-1, 1], [-1, -1]
    ]
    dest_dposes = [for dpos in mid_dposes : [dpos[0] * 2, dpos[1] * 2]]
    mids = [
      for dpos in mid_dposes
      : get_square_relative(square, dpos)
    ]
    valid_mid_idxs = [
      for i, square in mids : i
      if square == null ? false : is_occupied(square)
    ]
    dests = [
      for i in valid_mid_idxs
      : get_square_relative(square, dest_dposes[i])
    ]
    return = [
      for square in dests : square
      if square == null ? false : !is_occupied(square)
    ]
  }
}

// ===== HELPER FUNCTIONS ======================================
// Checks if square is occupied.
function "is_occupied" {
  params = [square]
  result = piece_at(square) != null
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

// ===== GAME STATE VALIDATORS ================================================
// Validators are called just after a move is taken. If any validator returns
// false, then the move is reversed - it cannot be completed.
//
// Validators receive 1 parameter - the last move and return true if the state
// is valid or false otherwise.

// No state validators in Halma.

// ===== INITIAL STATE ========================================================
// Initial state block specifies the initial placement of all the pieces.
initial_state {
  white_pieces = { for pos in starting_positions.white : pos => "piece" }
  black_pieces = { for pos in starting_positions.black : pos => "piece" }
}

// ===== VARIABLES =============================================
// Variables contain arbitrary data, which can be accessed in other blocks.
variables {
  starting_positions = {
    black = [
      "A8", "B8", "C8", "D8",
      "A7", "B7", "C7",
      "A6", "B6",
      "A5",
    ]
    white = [
      /*             */ "H4",
      /*       */ "G3", "H3",
      /* */ "F2", "G2", "H2",
      "E1", "F1", "G1", "H1",
    ]
  }
  max_moves_to_leave_start = 30
}

// ===== TURN =================================================================
// Turn block designates a function which controls the flow of a single turn
// (attribute "action"). Similarly to motion actions, turn action can interpret
// player choices via a choice generator. Read piece_types description for a
// more detailed description.
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
// Namely the function "pick_winner" and its helpers.

// This function is called at the start of every turn.
// Returns a tuple in form [is_finished, winner_color]. If is_finished == true
// and winner_color == null then draw is concluded.
function "pick_winner" {
  params = [game]
  result = (
    did_win("black")
    ? [true, "black"]
    : did_win("white")
    ? [true, "white"]
    : [false, null]
  )
}

// Checks if the given player meets the winning conditions.
composite_function "did_win" {
  params = [color]
  result = {
    opponent                = opponent_color(color)
    opponent_starting_poses = starting_positions[opponent]
    return = (
      is_all_occupied_by(opponent_starting_poses, color)
      || (
        count_moves_by(opponent) > max_moves_to_leave_start
        && is_any_occupied_by(opponent_starting_poses, opponent)
      )
    )
  }
}

// Checks if the given player occupies all of the given squares.
composite_function "is_all_occupied_by" {
  params = [squares, color]
  result = {
    pieces = [for square in squares : piece_at(square)]
    return = all([
      for piece in pieces : piece == null ? false : piece.color == color
    ]...)
  }
}

// Checks if the given player occupies any of the given squares.
composite_function "is_any_occupied_by" {
  params = [squares, color]
  result = {
    pieces = [for square in squares : piece_at(square)]
    return = any([
      for piece in pieces : piece == null ? false : piece.color == color
    ]...)
  }
}

// Returns number of moves done by the given player.
function "count_moves_by" {
  params = [color]
  result = sum([for move in game.record : 1 if move.piece.color == color]...)
}

// ===== ASSETS ===============================================================
// Assets are additional files, which can be used in rules. These files are
// copressed via gzip, and then ascii-encoded via base64. Any whitespace is
// ignored.
//
// Assets can be organized into an arbitrally nested tree.
assets = {
  piece_types = {
    "disk.svg" = <<EOF
      H4sIAOFLWmUAA1WRXW7DIBCE33MKRF8SqcbgJK1D40TqTZCNHVoMaCF2fPtunLo/aKVdfZod0HA8
      33pLBg3ReFdRwTgl2tW+Ma6r6DW1WUnPp9UxDh0ZjB7f/a2inHCy22NRgtsuVvSSUpB5Po4jG7fM
      Q5cXnPMct+hpRcixI62xtqKDgnWWBasmDVntrYfnp7ZtN5TEBP5TLwofgnfapUWDZrNmsijxQdUm
      TVK83U0zuFot9aCdb5oH+RU8XLPRNOkiBdsvwBqnaxUk+Ktr/sIPb9x/2pukwRpscrewRsWLAlCT
      dPjMhf5cS0kC5WLroa9orxKY21o8Ez7X98Be+HyKbbkT5f71sJmTwqy0tSZETWpMuigYhlxPOAl2
      oASQiRL7NPd8Djfv8IPuWZ9WX/hY10rPAQAA
      EOF
  }
}
