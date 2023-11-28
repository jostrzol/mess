// ===== BOARD ================================================================
// Board size definition.
board {
  height = 4
  width  = 3
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
  piece_type "lion" {
    presentation {
      white {
        icon = "/piece_types/lion.svg"
      }
      black {
        icon = "/piece_types/lion.svg"
        rotate = true
      }
    }
    motion {
      generator = "motion_neighbours"
    }
  }

  piece_type "hen" {
    presentation {
      white {
        icon = "/piece_types/hen.svg"
      }
      black {
        icon = "/piece_types/hen.svg"
        rotate = true
      }
    }
    motion {
      generator = "motion_hen"
    }
  }

  piece_type "giraffe" {
    presentation {
      white {
        icon = "/piece_types/giraffe.svg"
      }
      black {
        icon = "/piece_types/giraffe.svg"
        rotate = true
      }
    }
    motion {
      generator = "motion_neighbours_straight"
    }
  }

  piece_type "elephant" {
    presentation {
      white {
        icon = "/piece_types/elephant.svg"
      }
      black {
        icon = "/piece_types/elephant.svg"
        rotate = true
      }
    }
    motion {
      generator = "motion_neighbours_diagonal"
    }
  }

  piece_type "chick" {
    presentation {
      white {
        icon = "/piece_types/chick.svg"
      }
      black {
        icon = "/piece_types/chick.svg"
        rotate = true
      }
    }
    motion {
      generator = "motion_forward_straight"
      action    = "promote"
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

// ===== ACTION/CHOICE FUNCTIONS ==============================================
// Actions are executed after piece movement and can be parametrized with user
// decisions. To learn more, read description above piece_types.

// Exchanges a piece for a hen.
composite_function "promote" {
  params = [piece, src, dst, options]
  result = {
    owner     = owner_of(piece)
    forward_y = owner.forward_direction[1]
    last_y    = forward_y == 1 ? board.height - 1 : 0
    pos       = square_to_coords(dst)
    return    = cond_call(pos[1] == last_y, "place_new_piece", "hen", dst, owner.color)
  }
}

// ===== GAME STATE VALIDATORS ================================================
// Validators are called just after a move is taken. If any validator returns
// false, then the move is reversed - it cannot be completed.
//
// Validators receive 1 parameter - the last move and return true if the state
// is valid or false otherwise.

// No state validators in Dobutsu Shogi.

// ===== HELPER FUNCTIONS =====================================================
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

// ===== INITIAL STATE ========================================================
// Initial state block specifies the initial placement of all the pieces.
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

// ===== TURN =================================================================
// Turn block designates a function which controls the flow of a single turn
// (attribute "action"). Similarly to motion actions, turn action can interpret
// player choices via a choice generator. Read piece_types description for a
// more detailed description.
turn {
  choice_function = "turn_choices"
  action          = "turn"
}

function "turn_choices" {
  params = []
  result = {
    type = "unit"
    next_choices = [
      {
        message      = "Make a move"
        type         = "move"
        options      = []
        next_choices = []
      },
      {
        message = "Place a captured piece"
        type    = "piece_type"
        options = captured_piece_types(game.current_player)
        next_choices = [{
          message = "Choose an empty square"
          type    = "square"
          squares = empty_squares()
        }]
      },
    ]
  }
}

function "captured_piece_types" {
  params = [player]
  result = [for type, _ in player.captures : type]
}

composite_function "empty_squares" {
  params = []
  result = {
    squares = concat([for x in range(board.width) :
      [for y in range(board.height) : coords_to_square([x, y])]
    ]...)
    return = [for square in squares : square if piece_at(square) == null]
  }
}

composite_function "turn" {
  params = [options]
  result = {
    return = (options[1].type == "move"
      ? make_move(options[1].move, slice(options, 2, length(options)))
      : convert_and_release(
        game.current_player,
        options[1].piece_type.name,
        options[2].square
      )
    )
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

// ===== ASSETS ===============================================================
// Assets are additional files, which can be used in rules. These files are
// copressed via gzip, and then ascii-encoded via base64. Any whitespace is
// ignored.
//
// Assets can be organized into an arbitrally nested tree.
assets = {
  piece_types = {
    "lion.svg" = <<EOF
      H4sIAJmAXmUAA61V227aQBB9z1es3JdGSpbd2TuFPNTP/YG+UTDg1sHIkBD+vmfWBWJU1KpqhLJ7
      9jJzZs6Md7J7XYm352azmxbr/X47Ho0Oh4M8GNl2qxEppUY4UYhl3TTT4sNyuSzEbt+1Pyog7DI6
      NgDtdjav98ex/sRHH7uXphpXr9WmXSz6lcuB/v7joV7s1xfY1JtqPtuOu/Zls3i/+L2tN8PV53pf
      dU2NYWxPa4vZbj3rutlxvGk31Wn17LQQr3V1+Ny+TQsllLAOv+LpTuBvsurHPN/O9muxmBZfRJLB
      WGGCdCmJcgi1k0Z7YYzU1goy0quMkjE4arSMCdjKQCSMk9p7vhn0NSqvd42M2rCBkIQh/I+CkvQ2
      9mZDNIJAgb1hmec6DuflcMdKZbUgCy/ZvNHgq6QlyydJuhiE9tL3VgwAyOPGAJRXe7jmKCNOgobt
      KDRyYfgo+IZohbbSWuwqaZQRXkYfhqAcQiLpHa4pmZw+I43YOaUnTNIYl9MPwqBNGSjFgEMr2a8n
      hgqmQSoxNxgFU2xp0GejIfZH+Q6M2syXbAZORR6Uz0i5Xv0YezpaOJmCyVTUFSiHEMcdTAXpvWNg
      QRMoBQ4emAySGCXFJKL0JmaZFAmWDiohyKy7Q+UJ8ig8ywEp1BhyrJClISqvcJIuQHNkAnENyvfr
      uWv7Nhl3q28f1YPg333u1ryik3uATXoQMd0Xo992CafbM3Frw6/snyHo+ASWCeWWuMocDwHa5aMo
      KejOWqIbiXvG8H3uGY49Rt4LxCVNEfk3bCLk8JIL2SwkQhGhY1DSKb6rKdQLMgAQ2RqXhOZ7jnsh
      SMrZglEYJ5ZWnzziBvchkO6LN0bHR61MynKpKygK4wYaog0U2hyunc09y4VTsqJBJb6biBmnHviQ
      G4214Gj6GoAbNDoYKd61MrJ67HuYyL/Uiyz0IrIQTetbgrF5jhCp9Dm1A4zQAj4NQCnYU6CcP+cu
      XYhdnwIHrnmAoDQAf2bLn+hb/DR8Zh8xd9QQgwGZxBm07M2AdOYejfsHr1XT1NtdJeZ4FVBOKYL6
      /Dgt4NO5WIiOXwt8gQKmx/N06Ca/NmcJTg5vukG68UW5uEn/w8tkhCdswg/1091PVJSuodEHAAA=
      EOF
    "hen.svg" = <<EOF
      H4sIAJKAXmUAA61VXW/aQBB8z684uS+NlBz3fbcu5KF5bf9A3ygYcOtgZBwI/76zZ4XEqGkrNQhk
      z/k8Ozu7e0z3h7V4emi2+1mx6ftdOZkcj0d5tLLt1hOjlJpgRyFWddPMig+r1aoQ+75rf1ZAeMro
      1AC0u/mi7k+l/sRbb7vHpiqrQ7Vtl8th5WXD8P7tsV72mxfY1NtqMd+VXfu4Xb5e/NHW2/HqQ91X
      XVPjUrrnteV8v5l33fxUbttt9bx6DlqIQ10dP7dPs0IJJZzHt7i7EvhM18M13+/m/UYsZ8VXYYxU
      OghrZfJBfBFGSwrAUYboz2kPccpu/f2juhH8vc7pZhXF5PfMShofhQ3S6CTuLzDHIS2sk9EGYawM
      xvEzFeJ/RbVeWk9CJ2kTRx1jpEnKCO1xQXAjPSzSQXpDyB2YKPHekIjfHWF+NwARDIsXvN/+rplX
      jPc3yNzcCNLXb9vmhNYwyQymnREiawKCKngFA5OLQjuZy2ahIwhtYaDm90YYRY4oKnhCNCPOfxXu
      INwr/r2pG2ZGLQi0jgW8hqyVLMckSIY2EziN5CxL99LFyFiFnPIIG+TnQRMsjTnfUbmXCl4lqbwd
      4r9AhPceRVfSB8rSoAZ5JMUNYyAlEksNKUsfYRPzaEWpx5zvpxyyIs8SOhPDe3+B0ZvJe07BWZ/H
      j3IrKJe1ahlN5KYKjvgpew6B1vIjh6bGXCifhloaFbjVVBosIeJYmmlpaDRUE2z3PNA6Wp4wrylP
      kM7j7aALKFJGieIwXyEGJrS4ZBSZ0BvoTFmDhdU+95ODfMr0CUIjKBIjGoQFFZnAElcQ9bEclSQY
      SHKIiI3+wp8/1mHkeNU09W5fiQXOVQhLfEYtTrMCDlnnCtHxeQsfsdydzrdj8nxen4t7GWY6wfE8
      5T+hu6tfclvXaK0GAAA=
      EOF
    "giraffe.svg" = <<EOF
      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
      B8+nD/8H+/4NbwkSAAA=
      EOF
    "elephant.svg" = <<EOF
      H4sIAOd/XmUAA21U23LaMBB9z1do3JdmphHa1Z1CHspr+wN9o2DArYMZ44Tw9z0rFxIy9Xgsn71o
      L2el2fFlq16f2v1xXu2G4TCdTE6nkz5Z3fXbCRtjJrCo1KZp23n1abPZVOo49N2fGghaQecWoDss
      V81wntJXMX3on9t6Wr/U+269HiVvBqP/w6lZD7s32Db7erU8TPvueb9+L/zdNftb6VMz1H3bYJm6
      i2y9PO6Wfb88T/fdvr5Ir0Er9dLUp2/d67wyyijn8VaPdwrPbDuu5f+wHHZqPa9+qKzZW2W9dsmp
      xS0MmgIrDjrnAMDZKnbaRYJh0I6dYtLZiS5CSEH7zIoMLKMiq7ONsMRPYquINCdonY4QwyjHKC7Z
      F2RDsc2aIM46QopYBskk7EPKGu3FnbGIpY2aqMRkZwVxCuLuPQvyxkpyLvyzjajDapt4BKQY/hyU
      DTpl1JF0QHULwQY+iCZKp718IiJiCaUxxiWxgzQItCYry5pcFhlZ9wEtPmoJAVCO1cEHQS6Piw1i
      S6VYBg1UxJGcEJCCF2RDUOzRTjGVDhZlyOIQfck1je4plexo3NVoQ7Gw6kjamlJRxlz4fEOL0vMc
      3mmjG+uzfLUlncoIADNHwSiIPDKwAnyywnIGO5JylmSBnc2ja7zMgHQvf0Tfb0fw5/Xkyema9ttf
      nymaL2C6fPi+mvx3rEGvQX/BeS583WKpMnpZMrrGrK0vKFpJlg2oLtiDXIo6G1c80RggljPB2KAM
      LNrP4kE+CZDRooSlOAqpJIMqcwAx+dIol7O4MBKSUyKZgGlXUPKSgdVsC4zsy6hi2OBoZRpk23Jy
      oqdrd8aLoPQHjZH3vtxH5Zq46VHdts3hWKsV7giZO3aVWp3nFckIVaqHWKYRW/fn6+9tlHL1XPm4
      xLtEmU1w08zkPn28+wtkv7QneAUAAA==
      EOF
    "chick.svg" = <<EOF
      H4sIAIJ+XmUAA51Uy27bMBC85ysI9dIAyYrk8qnaOdTX9gd6Uy3ZVqtIhqxY8d93l6qd2GhRoIZg
      cpaznOVwpcXhuBWvz213WGa7cdwXeT5NE0wI/bDNtZQyJ0YmNk3bLrMPm80mE4dx6H/WhGiV0akl
      0O/LdTOeCvWJqY/DS1sX9bHu+qqaI2+EOf9xaqpx9wbbpqvX5b4Y+peueh/80TfddfS5GeuhbWgo
      zDlWlYddOQzlqej6rj5HL6KZODb19Ll/XWZSSGEsPdnTnaDfYjuPab4vx52oltlXoRCcF1qBN0as
      riEBGbRQAZSMQjnwnv+tjcTUEowNQlmQtKg1SExUF+yMhIpgMDCVoA6coaJiZJRLwehv0Op2FcFQ
      BVwP0rakFeY1qZnrwCVpTQNqQJ0WrXGMbERGSBkrxnOmJg4iqJAK0T4lOpOktFEzVSlKNaDpEDqC
      syhQQow6kehIqMB4ZLscBI2MoyW/DHgmWTo3XpABh7zvGSOlkJsKLDuFENHdoC/X1/Dt0n3zdRfD
      9vtH+SD4uU9dlyLahAfyOD4Ir++z/I+3TQKKiuH6nSEdsiB6LdDTtel/y3DP/W1nCZaOhW5ujd9z
      2t8Fz4ZHso2MURy3EKX9D7W6bZv9oRZram9FJauYifVpmXErSur9gdseog40PV2m1zLptbl4dhY8
      yyxyeksW/C14uvsFc2qdkjQEAAA=
      EOF
  }
}
