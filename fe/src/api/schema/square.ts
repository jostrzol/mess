import {Square} from "@/model/game/square"

export type SquareDto = [number, number]

export const squareToModel = (square: SquareDto): Square => {
  return square
}
