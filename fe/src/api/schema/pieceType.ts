import {PieceType} from "@/model/game/pieceType";

export interface PieceTypeDto {
  Name: string
}

const iconUriMap: Record<string, string> = {
  king: "/pieces/king.svg",
  queen: "/pieces/queen.svg",
  rook: "/pieces/rook.svg",
  bishop: "/pieces/bishop.svg",
  knight: "/pieces/knight.svg",
  pawn: "/pieces/pawn.svg",
};

export const pieceTypeToModel=(pieceType: PieceTypeDto): PieceType => {
  return {
    name: pieceType.Name,
    iconUri: iconUriMap[pieceType.Name],
  }
}
