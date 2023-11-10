import { GameState } from "@/model/game/gameState";
import { PieceDto, pieceToModel } from "./piece";
import {optionNodeToModel} from "./options";

export interface GameStateDto {
  TurnNumber: number;
  Pieces: PieceDto[];
  OptionTree: any;
  IsMyTurn: boolean
}

export const gameStateToModel = (state: GameStateDto): GameState => {
  return {
    turnNumber: state.TurnNumber,
    pieces: state.Pieces.map(pieceToModel),
    optionTree: optionNodeToModel(state.OptionTree),
    isMyTurn: state.IsMyTurn,
  };
};
