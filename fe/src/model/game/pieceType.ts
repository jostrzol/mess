import { Color } from "./color";

export interface PieceType {
  name: string;
  presentation: Record<Color, Presentation>;
}

export interface Presentation {
  symbol: string;
  icon?: string;
  rotate: boolean;
}
