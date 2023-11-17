import { Square } from "@/model/game/square";
import { useDroppable } from "@dnd-kit/core";
import clsx from "clsx";
import { HTMLAttributes } from "react";

export type DotType = "normal" | "danger";

export type TileProps = {
  square: Square;
  isDot?: boolean;
  dotType?: DotType;
  dotScale?: number;
} & HTMLAttributes<HTMLDivElement>;

export const Tile = ({
  square,
  isDot = false,
  dotType = "normal",
  dotScale = 1,
  children,
  ...props
}: TileProps) => {
  const { isOver, setNodeRef } = useDroppable({
    id: Square.toString(square),
    data: { square: square },
  });
  const dotColor = {
    normal: "bg-success-strong/90",
    danger: "bg-danger/90",
  }[dotType];
  const dotScaleEffective = !isDot ? 0 : isOver ? dotScale * 1.5 : dotScale;
  return (
    <div
      ref={setNodeRef}
      className={clsx(
        "min-h-[3rem]",
        "min-w-[3rem]",
        Square.isBlack(square) ? "bg-player-black" : "bg-player-white",
        "rounded-2xl",
        "relative",
      )}
      {...props}
    >
      <div>
        <div
          className={clsx(
            "absolute",
            "z-10",
            "top-1/2",
            "left-1/2",
            "w-1/4",
            "h-1/4",
            "rounded-full",
            dotColor,
            "transition-opacity",
            "pointer-events-none",
            "transition-transform",
          )}
          style={{
            transform: `translate(-50%, -50%) scale(${dotScaleEffective}) `,
          }}
        />
      </div>
      <div>{children}</div>
    </div>
  );
};
