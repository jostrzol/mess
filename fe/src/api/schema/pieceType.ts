import { Color } from "@/model/game/color";
import { PieceType, Presentation } from "@/model/game/pieceType";
import { ColorDto } from "./color";

export interface PieceTypeDto {
  Name: string;
  Presentation: Record<ColorDto, PresentationDto>;
}

export interface PresentationDto {
  Symbol: string;
  Icon?: string;
  Rotate: boolean;
}

export const pieceTypeToModel = (pieceType: PieceTypeDto): PieceType => {
  return {
    name: pieceType.Name,
    presentation: Object.fromEntries(
      Object.entries(pieceType.Presentation).map(([color, presentation]) => [
        color,
        presentationToModel(presentation),
      ]),
    ) as Record<Color, Presentation>,
  };
};

const presentationToModel = (presentation: PresentationDto): Presentation => {
  return {
    symbol: presentation.Symbol,
    icon: presentation.Icon,
    rotate: presentation.Rotate,
  };
};

export const pieceTypeToDto = (pieceType: PieceType): PieceTypeDto => {
  return {
    Name: pieceType.name,
  } as any;
};
