import { Board as BoardModel } from "@/model/game/board";
import { MoveMap } from "@/model/game/options";
import { Piece } from "@/model/game/piece";
import { Square } from "@/model/game/square";
import { useState } from "react";
import { Tile } from "./tile";

export type BoardProps = {
  pieces: Piece[];
  board: BoardModel;
  moveMap: MoveMap;
};

export const Board = ({ board, pieces, moveMap }: BoardProps) => {
  const [hoveredPiece, setHoveredPiece] = useState<Piece | null>(null);
  const gridTemplateColumns = `repeat(${board.width}, 1fr)`;
  const gridTemplateRows = `repeat(${board.width}, auto)`;
  const piecesMap = Object.fromEntries(
    pieces.map((piece) => [Square.toString(piece.square), piece]),
  );
  const destinations = (() => {
    if (hoveredPiece == null) {
      return [];
    }
    const from = Square.toString(hoveredPiece.square);
    return Object.keys(moveMap[from] ?? {});
  })();

  return (
    <div
      className="grid grid-flow-row max-h-full max-w-full h-fit aspect-square p-12"
      style={{ gridTemplateColumns, gridTemplateRows }}
    >
      {[...Array(board.height).keys()].flatMap((_, j) =>
        [...Array(board.width).keys()].map((_, i) => {
          const y = board.height - 1 - j;
          const x = i;
          const square = Square.toString([x, y]);
          return (
            <Tile
              key={square}
              color={(x + y) % 2 == 1 ? "white" : "black"}
              piece={piecesMap[square]}
              canMove={square in moveMap}
              onPieceHovered={(piece) => setHoveredPiece(piece)}
              onPieceUnhovered={(piece) =>
                hoveredPiece == piece && setHoveredPiece(null)
              }
              isMoveProjected={destinations.includes(square)}
            />
          );
        }),
      )}
    </div>
  );
};
