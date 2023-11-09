import { GameState } from "@/model/game/gameState";
import { PieceDto, pieceToModel } from "./piece";
import {optionNodeToModel} from "./options";

export interface GameStateDto {
  Pieces: PieceDto[];
  OptionTree: any;
}

export const gameStateToModel = (state: GameStateDto): GameState => {
  return {
    pieces: state.Pieces.map(pieceToModel),
    optionTree: optionNodeToModel(state.OptionTree),
  };
};
