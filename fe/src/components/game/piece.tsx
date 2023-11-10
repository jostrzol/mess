import { useBoard } from "@/contexts/boardContext";
import { useOptions } from "@/contexts/optionContext";
import * as model from "@/model/game/piece";
import { Square } from "@/model/game/square";
import { useDraggable } from "@dnd-kit/core";
import { CSS } from "@dnd-kit/utilities";
import clsx from "clsx";
import { ReactSVG } from "react-svg";

export interface PieceProps {
  piece: model.Piece;
}

export const Piece = ({ piece }: PieceProps) => {
  const { moveMap } = useOptions();
  const { hoveredSquare } = useBoard();

  const moves = moveMap[Square.toString(piece.square)];
  const canMove = moves !== undefined
  const canDrop = hoveredSquare && Square.toString(hoveredSquare) in (moves ?? {});

  const { attributes, listeners, setNodeRef, transform, isDragging } =
    useDraggable({
      id: Square.toString(piece.square),
      disabled: !canMove,
      data: { piece: piece },
    });

  const style = { transform: CSS.Translate.toString(transform) };
  return (
    <div
      ref={setNodeRef}
      className={clsx(
        "relative",
        canMove && "hover:scale-110",
        !canMove && "cursor-default",
        isDragging && ["z-20", "scale-110", "cursor-none"],
        isDragging && !canDrop && ["opacity-50"],
      )}
      style={style}
      {...listeners}
      {...attributes}
    >
      <ReactSVG
        className={clsx(
          piece.color == "white" ? "player-white" : "player-black",
        )}
        src={piece.type.iconUri}
      />
    </div>
  );
};
