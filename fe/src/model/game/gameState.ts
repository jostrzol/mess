import { OptionNode } from "./options";
import { Piece } from "./piece";

export interface GameState {
  turnNumber: number;
  pieces: Piece[];
  optionTree: OptionNode;
  isMyTurn: boolean;
}
