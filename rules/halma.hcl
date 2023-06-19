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
  piece_type "piece" {
    motion {
      generator = "motion_neighbours"
    }
    motion {
      generator = "motion_jump"
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

// ===== INITIAL STATE =========================================
initial_state {
  white_pieces = { for pos in white_starting_positions : pos => "piece" }
  black_pieces = { for pos in black_starting_positions : pos => "piece" }
}

// ===== VARIABLES =============================================
variables {
  black_starting_positions = [
    "A8", "B8", "C8", "D8",
    "A7", "B7", "C7",
    "A6", "B6",
    "A5",
  ]
  white_starting_positions = [
    /*             */ "H4",
    /*       */ "G3", "H3",
    /* */ "F2", "G2", "H2",
    "E1", "F1", "G1", "H1",
  ]
}

// ===== GAME RESOLVING FUNCTIONS ==============================
// Namely the function "pick_winner" and its helpers

// This function is called at the start of every turn.
// Returns a tuple in form [is_finished, winner_color]. If is_finished == true
// and winner_color == null then draw is concluded.
function "pick_winner" {
  params = [game]
  result = (
    is_all_occupied_by(white_starting_positions, "black")
    ? [true, "black"]
    : (is_all_occupied_by(black_starting_positions, "white")
      ? [true, "white"]
      : [false, null]
    )
  )
}

// Checks if the given player occupies all of the given squares
composite_function "is_all_occupied_by" {
  params = [squares, color]
  result = {
    pieces = [for square in squares : piece_at(square)]
    return = all([
      for piece in pieces : piece == null ? false : piece.color == color
    ]...)
  }
}
