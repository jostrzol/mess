import { Piece } from "./piece";

export interface GameState {
  turnNumber: number;
  pieces: Piece[];
  isMyTurn: boolean;
}
