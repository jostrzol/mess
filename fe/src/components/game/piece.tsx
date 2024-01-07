import { useBoard } from "@/contexts/boardContext";
import { useGameState } from "@/contexts/gameStateContext";
import { useOptions } from "@/contexts/optionContext";
import { useStaticData } from "@/contexts/staticDataContext";
import { Color } from "@/model/game/color";
import * as model from "@/model/game/piece";
import { Presentation } from "@/model/game/pieceType";
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
  const presentation = piece.type.presentation[piece.color];

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
        isDragging ? "z-30" : "z-10",
        isDragging && !canDrop && ["opacity-50"],
      )}
      style={style}
      {...listeners}
      {...attributes}
    >
      <div
        className={clsx(
          "transition-transform",
          (canMove || isDragging) && "hover:scale-125",
        )}
      >
        <PieceIcon
          className="p-1"
          presentation={presentation}
          color={piece.color}
        />
      </div>
    </div>
  );
};

export const PieceIcon = ({
  presentation,
  blockRotation = false,
  color,
  className,
}: {
  presentation: Presentation;
  blockRotation?: boolean;
  color: Color;
  className?: string;
}) => {
  const { assetUrl } = useStaticData();

  const icon =
    presentation.icon === undefined ? (
      <svg
        viewBox="0 0 100 100"
        className={clsx(
          color == "white" ? "player-white" : "player-black",
          "text-player",
          className,
        )}
        style={{
          font: "bold 80px Century Gothic, Arial",
          fill: "var(--player-color)",
          stroke: "var(--opponent-color)",
          strokeWidth: 4,
          strokeLinejoin: "round",
          strokeLinecap: "round",
        }}
      >
        <text
          y="50%"
          x="50%"
          textAnchor="middle"
          dominantBaseline="central"
          style={{ fontSize: "80px" }}
        >
          {presentation.symbol}
        </text>
      </svg>
    ) : (
      <ReactSVG
        className={clsx(
          color == "white" ? "player-white" : "player-black",
          "transition-transform",
          className,
        )}
        src={assetUrl(presentation.icon)}
      />
    );

  return (
    <div
      className={clsx(presentation.rotate && !blockRotation && "rotate-180")}
    >
      {icon}
    </div>
  );
};
