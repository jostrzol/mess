import { Color } from "./color";
import { PieceType } from "./pieceType";
import { Square } from "./square";

export interface Piece {
  type: PieceType;
  color: Color;
  square: Square;
}
