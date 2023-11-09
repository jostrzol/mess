import { Color } from "@/model/game/color";
import { Piece } from "@/model/game/piece";
import clsx from "clsx";
import * as component from "./piece";

export type TileProps = {
  color: Color;
  piece?: Piece;
  canMove: boolean;
  onPieceHovered?: (piece: Piece) => any;
  onPieceUnhovered?: (piece: Piece) => any;
  isMoveProjected?: boolean;
};

export const Tile = ({
  color,
  piece,
  canMove,
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
        piece && canMove && "hover:cursor-pointer",
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
          piece ? "bg-danger/90" : "bg-success-strong/90",
          "transition-opacity",
          isMoveProjected || "opacity-0",
          "pointer-events-none",
        )}
      />
      <div
        className={clsx(canMove && "hover:scale-110", "transition-transform")}
        onPointerEnter={() => piece && onPieceHovered?.(piece)}
        onPointerLeave={() => piece && onPieceUnhovered?.(piece)}
      >
        {piece && <component.Piece piece={piece} />}
      </div>
      {/* Needed to make the parent div expand.
      Coulnd't get it to work without the image */}
      <svg className="invisible">
        <rect width="1" height="1" />
      </svg>
    </div>
  );
};
