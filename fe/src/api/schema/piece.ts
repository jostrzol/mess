import { Piece } from "@/model/game/piece";
import { ColorDto, colorToModel } from "./color";
import { PieceTypeDto, pieceTypeToModel } from "./pieceType";
import { SquareDto, squareToModel } from "./square";

export interface PieceDto {
  Type: PieceTypeDto;
  Color: ColorDto;
  Square: SquareDto;
}

export const pieceToModel = (piece: PieceDto): Piece => ({
  type: pieceTypeToModel(piece.Type),
  square: squareToModel(piece.Square),
  color: colorToModel(piece.Color),
});
