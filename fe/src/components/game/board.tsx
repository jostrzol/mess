import { BoardProvider, useBoard } from "@/contexts/boardContext";
import { useGameState } from "@/contexts/gameStateContext";
import { useOptions } from "@/contexts/optionContext";
import { Board as BoardModel } from "@/model/game/board";
import { Move } from "@/model/game/move";
import { Piece as PieceModel } from "@/model/game/piece";
import { Square } from "@/model/game/square";
import { DndContext } from "@dnd-kit/core";
import clsx from "clsx";
import { Piece } from "./piece";
import { Tile } from "./tile";

export type BoardProps = {
  board: BoardModel;
};

export const Board = (props: BoardProps) => {
  return (
    <BoardProvider>
      <BoardWrapped {...props} />
    </BoardProvider>
  );
};

const BoardWrapped = ({ board }: BoardProps) => {
  const gridTemplateColumns = `repeat(${board.width}, 1fr)`;
  const gridTemplateRows = `repeat(${board.width}, auto)`;

  const { pieceMap } = useGameState();
  const { choose } = useOptions();
  const { dispatch, destinations, draggedPiece } = useBoard();

  return (
    <DndContext
      modifiers={[]}
      onDragStart={(e) => {
        const piece: PieceModel = e.active.data.current!.piece;
        dispatch({ type: "Dragged", piece: piece });
      }}
      onDragOver={(e) => {
        if (e.over === null) {
          dispatch({ type: "Unhovered" });
          return;
        }

        const over: Square = e.over.data.current!.square;
        dispatch({ type: "Hovered", square: over });
      }}
      onDragCancel={(e) => {
        const piece: PieceModel = e.active.data.current!.piece;
        dispatch({ type: "Dropped", piece: piece });
      }}
      onDragEnd={(e) => {
        const piece: PieceModel = e.active.data.current!.piece;
        dispatch({ type: "Dropped", piece: piece });

        if (e.over === null) return;

        const destination: Square = e.over.data.current!.square;
        const move: Move = {
          from: piece.square,
          to: destination,
        };

        console.log(`move: ${JSON.stringify(move)}`);
      }}
    >
      <div
        className={clsx(
          "p-12",
          "w-full",
          "h-full",
          draggedPiece && ["cursor-none", "[&_*]:cursor-none"],
        )}
      >
        <div
          className={clsx(
            "grid",
            "grid-flow-row",
            "max-h-full",
            "max-w-full",
            "h-fit",
            "aspect-square",
            "m-auto",
          )}
          style={{ gridTemplateColumns, gridTemplateRows }}
        >
          {BoardModel.MapSquares(board, (square, key) => {
            const piece = pieceMap[key];
            return (
              <Tile
                key={key}
                square={square}
                isDot={destinations.includes(key)}
                dotType={piece ? "danger" : "normal"}
              >
                {piece && <Piece piece={piece} />}
              </Tile>
            );
          })}
        </div>
      </div>
    </DndContext>
  );
};
