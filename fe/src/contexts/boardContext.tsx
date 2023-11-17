import { Piece } from "@/model/game/piece";
import { Square } from "@/model/game/square";
import {
  Dispatch,
  ReactNode,
  createContext,
  useContext,
  useReducer,
} from "react";
import { useGameState } from "./gameStateContext";
import { useOptions } from "./optionContext";

export const BoardContext = createContext<BoardContextValue>(null!);

export const useBoard = () => {
  return useContext(BoardContext);
};

interface Hovered {
  type: "Hovered";
  square: Square;
}

interface Unhovered {
  type: "Unhovered";
}

interface Dragged {
  type: "Dragged";
  piece: Piece;
}

interface Dropped {
  type: "Dropped";
  piece: Piece;
}

interface Dropped {
  type: "Dropped";
  piece: Piece;
}

export type Action = Hovered | Unhovered | Dragged | Dropped;

export interface BoardContextValue {
  hoveredSquare?: Square;
  draggedPiece?: Piece;
  focusedPiece?: Piece;
  destinations: string[];
  dispatch: Dispatch<Action>;
}

interface State {
  hoveredSquare?: Square;
  draggedPiece?: Piece;
}

export const BoardProvider = ({ children }: { children?: ReactNode }) => {
  const [state, dispatch] = useReducer(reducer, {});
  const { pieceMap } = useGameState();
  const { moveMap } = useOptions();

  const { hoveredSquare, draggedPiece } = state;
  const hoveredPiece =
    hoveredSquare && pieceMap[Square.toString(hoveredSquare)];
  const focusedPiece = draggedPiece || hoveredPiece;
  const destinations = (() => {
    if (focusedPiece == null) {
      return [];
    }
    const from = Square.toString(focusedPiece.square);
    return Object.keys(moveMap[from] ?? {});
  })();

  return (
    <BoardContext.Provider
      value={{
        hoveredSquare,
        draggedPiece,
        focusedPiece,
        destinations,
        dispatch,
      }}
    >
      {children}
    </BoardContext.Provider>
  );
};

const reducer = (boardState: State, action: Action) => {
  switch (action.type) {
    case "Hovered":
      return { ...boardState, hoveredSquare: action.square };
    case "Unhovered":
      return { ...boardState, hoveredSquare: undefined };
    case "Dragged":
      return { ...boardState, draggedPiece: action.piece };
    case "Dropped":
      return { ...boardState, draggedPiece: undefined };
    default:
      return boardState;
  }
};
