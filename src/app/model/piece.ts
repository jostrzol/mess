import {PieceType} from "./pieceType"

export interface Piece {
  type: PieceType
  location: [number, number]
}
