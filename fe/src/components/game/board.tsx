import { BoardProvider, useBoard } from "@/contexts/boardContext";
import { useGameState } from "@/contexts/gameStateContext";
import { useOptions } from "@/contexts/optionContext";
import { Board as BoardModel } from "@/model/game/board";
import { Piece as PieceModel } from "@/model/game/piece";
import { Square } from "@/model/game/square";
import { DndContext } from "@dnd-kit/core";
import clsx from "clsx";
import {Piece} from "./piece";
import {Tile} from "./tile";

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
  const gridTemplateRows = `repeat(${board.height}, 1fr)`;

  const { pieceMap, isMyTurn } = useGameState();
  const { choose, moveMap } = useOptions();
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
        const pieceMoves = moveMap[Square.toString(piece.square)] ?? {};
        const options = pieceMoves[Square.toString(destination)] ?? [];
        // TODO: choose option if options.length > 1
        const option = options.pop();

        if (option) {
          choose(option.node, option.datum);
        }
      }}
    >
      <div
        className={clsx(
          "grow",
          "portrait:w-11/12",
          "aspect-square",
          "flex",
          "flex-col",
          "justify-center",
          draggedPiece && ["cursor-none", "[&_*]:cursor-none"],
        )}
      >
        <div
          className={clsx(
          "p-4",
          "grid",
          "grid-flow-row",
          )}
          style={{ gridTemplateColumns, gridTemplateRows }}
          onPointerLeave={() => !draggedPiece && dispatch({"type": "Unhovered"})}
        >
          {BoardModel.MapSquares(board, (square, key) => {
          const piece = pieceMap[key];
          return (
            <Tile
            key={key}
              square={square}
              isDot={destinations.includes(key)}
              dotType={piece ? "danger" : "normal"}
              dotScale={isMyTurn ? 1 : 0.6}
              onPointerOver={() => !draggedPiece && dispatch({"type": "Hovered", square: square})}
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
