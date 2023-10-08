import { Color } from "./color";
import { PieceType } from "./pieceType";

export interface Piece {
  type: PieceType;
  color: Color;
  location: [number, number];
}
