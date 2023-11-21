import { BoardProvider, useBoard } from "@/contexts/boardContext";
import { useGameState } from "@/contexts/gameStateContext";
import { useOptions } from "@/contexts/optionContext";
import { Board as BoardModel } from "@/model/game/board";
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
  const gridTemplateRows = `repeat(${board.height}, 1fr)`;

  const { pieceMap, isMyTurn } = useGameState();
  const { choose, moveMap, squareMap } = useOptions();
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
        const routeItem = pieceMoves[Square.toString(destination)];

        if (routeItem === undefined) return

        choose(routeItem);
      }}
    >
      <div
        className={clsx(
          "grow",
          "portrait:w-10/12",
          "flex",
          "flex-col",
          "justify-center",
          draggedPiece && ["cursor-none", "[&_*]:cursor-none"],
        )}
        style={{aspectRatio: `${board.width} / ${board.height}`}}
      >
        <div
          className={clsx("grid", "grid-flow-row")}
          style={{ gridTemplateColumns, gridTemplateRows }}
          onPointerLeave={() =>
            !draggedPiece && dispatch({ type: "Unhovered" })
          }
        >
          {BoardModel.MapSquares(board, (square, key) => {
            const piece = pieceMap[key];
            const squareRouteItem = squareMap[key]
            return (
              <Tile
                key={key}
                square={square}
                isDot={destinations.includes(key)}
                dotType={piece ? "danger" : "normal"}
                dotScale={isMyTurn ? 1 : 0.6}
                isRing={squareRouteItem !== undefined}
                ringScale={isMyTurn ? 1 : 0.6}
                onPointerOver={() =>
                  !draggedPiece && dispatch({ type: "Hovered", square: square })
                }
                onClick={() => squareRouteItem && choose(squareRouteItem)}
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
