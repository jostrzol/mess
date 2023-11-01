import { Board as BoardModel } from "@/model/board";
import { Piece as PieceModel } from "@/model/piece";
import { Square } from "@/model/square";
import { useState } from "react";
import { Tile } from "./tile";

export type BoardProps = {
  pieces: PieceModel[];
  board: BoardModel;
};

export const Board = ({ board, pieces }: BoardProps) => {
  const [hoveredPiece, setHoveredPiece] = useState<PieceModel | null>(null);
  const gridTemplateColumns = `repeat(${board.width}, 1fr)`;
  const gridTemplateRows = `repeat(${board.width}, auto)`;
  const piecesMap = Object.fromEntries(
    pieces.map((piece) => [Square.toString(piece.square), piece]),
  );
  return (
    <div
      className="grid grid-flow-row max-h-full max-w-full h-fit aspect-square p-12"
      style={{ gridTemplateColumns, gridTemplateRows }}
    >
      {[...Array(board.height).keys()].flatMap((_, j) =>
        [...Array(board.width).keys()].map((_, i) => {
          const square = Square.toString([i, j]);
          return (
            <Tile
              key={square}
              color={(i + j) % 2 == 0 ? "white" : "black"}
              piece={piecesMap[square]}
              onPieceHovered={(piece) => setHoveredPiece(piece)}
              onPieceUnhovered={(piece) =>
                hoveredPiece == piece && setHoveredPiece(null)
              }
              isMoveProjected={
                (hoveredPiece || false) &&
                PieceModel.hasValidMove(hoveredPiece, [i, j])
              }
            />
          );
        }),
      )}
    </div>
  );
};
