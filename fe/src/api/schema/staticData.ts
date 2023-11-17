import { Board } from "@/model/game/board";
import { StaticData } from "@/model/game/gameStaticData";
import { UUID } from "crypto";
import { ColorDto, colorToModel } from "./color";

export interface StaticDataDto {
  ID: UUID;
  BoardSize: BoardSizeDto;
  MyColor: ColorDto;
}

export interface BoardSizeDto {
  Width: number;
  Height: number;
}

export const staticDataToModel = (staticData: StaticDataDto): StaticData => ({
  id: staticData.ID,
  board: boardSizeToModel(staticData.BoardSize),
  myColor: colorToModel(staticData.MyColor),
});

const boardSizeToModel = (boardSize: BoardSizeDto): Board => ({
  height: boardSize.Height,
  width: boardSize.Width,
});
