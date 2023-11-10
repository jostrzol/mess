import {Move} from "@/model/game/move";
import {SquareDto} from "./square";

export interface MoveDto {
  From: SquareDto
  To: SquareDto
}

export const moveToModel = (move: MoveDto): Move =>{
  return {
    from: move.From,
    to: move.To
  }
}

export const moveToDto = (move: Move): MoveDto =>{
  return {
    From: move.from,
    To: move.to
  }
}
