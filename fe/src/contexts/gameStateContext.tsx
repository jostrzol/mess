import { GameState } from "@/model/game/gameState";
import { Piece } from "@/model/game/piece";
import { Square } from "@/model/game/square";
import { ReactNode, createContext, useContext } from "react";

export const GameStateContext = createContext<GameStateContextValue>(null!);

export const useGameState = () => {
  return useContext(GameStateContext);
}

export interface GameStateContextValue extends GameState {
  pieceMap: Record<string, Piece>;
}

export const GameStateProvider = ({
  state,
  children,
}: {
  state: GameState;
  children?: ReactNode;
}) => {
  const pieceMap = Object.fromEntries(
    state?.pieces.map((piece) => [Square.toString(piece.square), piece]) ?? [],
  );
  return (
    <GameStateContext.Provider value={{ pieceMap, ...state }}>
      {children}
    </GameStateContext.Provider>
  );
};
