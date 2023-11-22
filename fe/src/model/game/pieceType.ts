import { Color } from "./color";

export interface PieceType {
  name: string;
  representation: Record<Color, Representation>;
}

export interface Representation {
  symbol: string;
  icon?: string;
}
