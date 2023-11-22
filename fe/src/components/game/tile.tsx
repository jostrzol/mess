import { Square } from "@/model/game/square";
import { useDroppable } from "@dnd-kit/core";
import clsx from "clsx";
import { HTMLAttributes } from "react";

export type DotType = "normal" | "danger";

export type TileProps = {
  square: Square;
  isRing?: boolean;
  ringScale?: number;
  isDot?: boolean;
  dotType?: DotType;
  dotScale?: number;
} & HTMLAttributes<HTMLDivElement>;

export const Tile = ({
  square,
  isDot = false,
  dotType = "normal",
  dotScale = 1,
  isRing = false,
  ringScale = 1,
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
  const ringScaleEffective = !isRing ? 0 : ringScale;

  return (
    <div
      ref={setNodeRef}
      className={clsx(
        "min-h-[3rem]",
        "min-w-[3rem]",
        "aspect-square",
        Square.isBlack(square) ? "bg-player-black" : "bg-player-white",
        "rounded-2xl",
        "relative",
      )}
      {...props}
    >
      <div>
        <div
          className={clsx(
            "absolute top-1/2 left-1/2 z-20",
            "w-1/4 h-1/4",
            "rounded-full",
            "transition-transform",
            "pointer-events-none",
            dotColor,
          )}
          style={{
            transform: `translate(-50%, -50%) scale(${dotScaleEffective}) `,
          }}
        />
        <div
          className={clsx(
            "absolute top-1/2 left-1/2 z-20",
            "w-1/4 h-1/4",
            "transition-transform",
          )}
          style={{
            transform: `translate(-50%, -50%) scale(${ringScaleEffective}) `,
          }}
        >
          <div
            className={clsx(
              "w-full h-full",
              "rounded-full",
              "border-8 border-success-strong/90",
              "transition-transform hover:scale-125",
              "cursor-pointer",
            )}
          />
        </div>
      </div>
      <div className="w-full h-full">{children}</div>
    </div>
  );
};
