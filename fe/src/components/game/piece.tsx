import * as model from "@/model/game/piece";
import clsx from "clsx";
import { ReactSVG } from "react-svg";

export interface PieceProps {
  piece: model.Piece;
}

export const Piece = ({ piece }: PieceProps) => {
  return (
    <ReactSVG
      className={clsx(piece.color == "white" ? "player-white" : "player-black")}
      src={piece.type.iconUri}
    />
  );
};
