import { Color } from "@/model/game/color";

export type ColorDto = "black" | "white";

export const colorToModel = (color: ColorDto): Color => {
  return color;
};
