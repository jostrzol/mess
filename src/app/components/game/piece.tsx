import { Piece as PieceConfig } from "@/app/model/piece";
import clsx from "clsx";
import { ReactSVG } from "react-svg";

export interface PieceProps {
  piece: PieceConfig;
}

export const Piece = ({ piece }: PieceProps) => {
  return (
    <ReactSVG
      className={clsx(
        piece.color == "white" ? "player-white" : "player-black",
        "hover:scale-110",
        "transition-transform",
      )}
      src={piece.type.iconUri}
    />
  );
};
