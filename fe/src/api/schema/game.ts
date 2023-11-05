import {Color} from "@/model/color"
import {GameState} from "@/model/gameState"
import {Piece} from "@/model/piece"

export interface GameStateDto {
  Pieces: {
    Type: { Name: string }
    Color: string
    Square: [number, number]
  }[]
}

const iconUriMap: Record<string, string> = {
    king:  "/pieces/king.svg",
    queen:  "/pieces/queen.svg",
    rook:  "/pieces/rook.svg",
    bishop:  "/pieces/bishop.svg",
    knight:  "/pieces/knight.svg",
    pawn:  "/pieces/pawn.svg",
}

export const gameStateToModel = (state: GameStateDto): GameState => {
  return {
    pieces: state.Pieces.map((piece): Piece => ({
      type: {
        name: piece.Type.Name,
        iconUri: iconUriMap[piece.Type.Name],
      },
      square: piece.Square,
      color: piece.Color as Color,
      validMoves: []
    }) )
  };
};
