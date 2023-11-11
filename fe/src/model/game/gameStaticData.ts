import { UUID } from "crypto";
import { Board } from "./board";
import { Color } from "./color";

export interface StaticData {
  id: UUID;
  board: Board;
  myColor: Color;
}
