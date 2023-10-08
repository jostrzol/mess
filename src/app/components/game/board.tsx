import { Board as BoardModel } from "@/app/model/board";
import { Color } from "@/app/model/color";
import { Piece as PieceModel } from "@/app/model/piece";
import { Square } from "@/app/model/square";
import clsx from "clsx";
import { useState } from "react";
import { Piece } from "./piece";

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

type TileProps = {
  color: Color;
  piece?: PieceModel;
  onPieceHovered?: (piece: PieceModel) => any;
  onPieceUnhovered?: (piece: PieceModel) => any;
  isMoveProjected?: boolean;
};

const Tile = ({
  color,
  piece,
  onPieceHovered,
  onPieceUnhovered,
  isMoveProjected = false,
}: TileProps) => {
  return (
    <div
      className={clsx(
        "min-h-[3rem]",
        "min-w-[3rem]",
        color == "white" ? "bg-player-white" : "bg-player-black",
        "rounded-2xl",
        "relative",
        piece && "hover:cursor-pointer",
      )}
      style={{
        overflowClipMargin: "content-box",
        overflow: "clip",
      }}
    >
      <div
        className={clsx(
          "absolute",
          "z-10",
          "top-1/2",
          "left-1/2",
          "-translate-x-1/2",
          "-translate-y-1/2",
          "w-4",
          "h-4",
          "rounded-full",
          piece ? "bg-danger" : "bg-primary/80",
          "transition-opacity",
          isMoveProjected || "opacity-0",
          "pointer-events-none",
        )}
      />
      <div
        className={clsx("hover:scale-110", "transition-transform")}
        onPointerEnter={() => piece && onPieceHovered?.(piece)}
        onPointerLeave={() => piece && onPieceUnhovered?.(piece)}
      >
        {piece && <Piece piece={piece} />}
      </div>
      {/* Needed to make the parent div expand.
      Coulnd't get it to work without the image */}
      <svg className="invisible">
        <rect width="1" height="1" />
      </svg>
    </div>
  );
};
