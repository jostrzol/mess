import {OptionNode} from "./options";
import {Piece} from "./piece";

export interface GameState {
  pieces: Piece[];
  optionTree: OptionNode
}
