import { PieceType, Representation } from "@/model/game/pieceType";
import { ColorDto } from "./color";
import {Color} from "@/model/game/color";

export interface PieceTypeDto {
  Name: string;
  Representation: Record<ColorDto, RepresentationDto>;
}

export interface RepresentationDto {
  Symbol: string;
  Icon: string;
}

export const pieceTypeToModel = (pieceType: PieceTypeDto): PieceType => {
  return {
    name: pieceType.Name,
    representation: Object.fromEntries(
      Object.entries(pieceType.Representation).map(
        ([color, repr]) => [color, representationToModel(repr)],
      ),
    ) as Record<Color, Representation>,
  };
};

const representationToModel = (
  representation: RepresentationDto,
): Representation => {
  return {
    symbol: representation.Symbol,
    icon: representation.Icon,
  };
};

export const pieceTypeToDto = (pieceType: PieceType): PieceTypeDto => {
  return {
    Name: pieceType.name,
  } as any;
};
