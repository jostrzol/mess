import { useBoard } from "@/contexts/boardContext";
import { useGameState } from "@/contexts/gameStateContext";
import { useOptions } from "@/contexts/optionContext";
import { useStaticData } from "@/contexts/staticDataContext";
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
  const { myColor } = useStaticData();
  const { moveMap, isReady } = useOptions();
  const { hoveredSquare } = useBoard();
  const { isMyTurn } = useGameState();

  const isMine = piece.color === myColor;
  const moves = moveMap[Square.toString(piece.square)];
  const canMove = isMyTurn && moves !== undefined;
  const canDrop =
    hoveredSquare && Square.toString(hoveredSquare) in (moves ?? {});

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
        !isMine
          ? "cursor-not-allowed"
          : isMyTurn && !isReady
          ? "cursor-wait"
          : isDragging
          ? "cursor-none"
          : null,
        !canMove && "cursor-default",
        isDragging && ["z-20"],
        isDragging && !canDrop && ["opacity-50"],
      )}
      style={style}
      {...listeners}
      {...attributes}
    >
      <ReactSVG
        className={clsx(
          piece.color == "white" ? "player-white" : "player-black",
          "transition-transform",
          (canMove || isDragging) && "hover:scale-110",
        )}
        src={piece.type.iconUri}
      />
    </div>
  );
};
