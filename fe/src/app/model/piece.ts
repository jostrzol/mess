import { Color } from "./color";
import { PieceType } from "./pieceType";
import { Square } from "./square";

export interface Piece {
  type: PieceType;
  color: Color;
  square: Square;
  validMoves: Square[];
}

export namespace Piece {
  export const hasValidMove = (piece: Piece, move: Square): boolean => {
    const found = piece.validMoves.find((validMove) => Square.equals(validMove, move));
    return found !== undefined;
  };
}
