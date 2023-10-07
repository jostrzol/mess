import { Board as BoardConfig } from "@/app/model/board";
import { Color } from "@/app/model/color";
import { Piece } from "@/app/model/piece";
import { locationToString } from "@/app/utils/functions";
import clsx from "clsx";

export type BoardProps = {
  pieces: Piece[];
  board: BoardConfig;
};

export const Board = ({ board, pieces }: BoardProps) => {
  const gridTemplateColumns = `repeat(${board.width}, 1fr)`;
  const gridTemplateRows = `repeat(${board.width}, auto)`;
  const piecesMap = Object.fromEntries(
    pieces.map((piece) => [locationToString(piece.location), piece]),
  );
  return (
    <div
      className="grid grid-flow-row max-h-full max-w-full h-fit aspect-square p-12"
      style={{ gridTemplateColumns, gridTemplateRows }}
    >
      {[...Array(board.height).keys()].flatMap((_, i) => {
        const x = [...Array(board.width).keys()].map((_, j) => {
          const location = locationToString([i, j]);
          return (
            <Tile
              key={location}
              color={(i + j) % 2 == 0 ? "white" : "black"}
              piece={piecesMap[location]}
            />
          );
        });
        return x;
      })}
    </div>
  );
};

type TileProps = {
  color: Color;
  piece?: Piece;
};

const Tile = ({ color, piece }: TileProps) => {
  return (
    <div
      className={clsx(
        "min-h-[3rem]",
        "min-w-[3rem]",
        color == "white" ? "bg-player-white" : "bg-player-black",
        "rounded-2xl",
      )}
      style={{
        overflowClipMargin: "content-box",
        overflow: "clip",
      }}
    >
      {/* Needed to make the parent div expand.
      Coulnd't get it to work without the image */}
      <svg className="invisible">
        <rect width="1" height="1" />
      </svg>
    </div>
  );
};
