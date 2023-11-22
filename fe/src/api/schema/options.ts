import { OptionNode, Route, RouteItem, Unit } from "@/model/game/options";
import { Square } from "@/model/game/square";
import { MoveDto, moveToDto, moveToModel } from "./move";
import { PieceTypeDto, pieceTypeToDto, pieceTypeToModel } from "./pieceType";
import { SquareDto, squareToDto, squareToModel } from "./square";

export type OptionNodeDto =
  | PieceTypeOptionNodeDto
  | SquareOptionNodeDto
  | MoveOptionNodeDto
  | UnitOptionNodeDto;

export type OptionDto = PieceTypeDto | SquareDto | MoveDto | UnitDto;

export type UnitDto = {};

export interface PieceTypeOptionNodeDto
  extends BaseOptionNodeDto<PieceTypeDto> {
  Type: "PieceType";
}

export interface SquareOptionNodeDto extends BaseOptionNodeDto<SquareDto> {
  Type: "Square";
}

export interface MoveOptionNodeDto extends BaseOptionNodeDto<MoveDto> {
  Type: "Move";
}

export interface UnitOptionNodeDto extends BaseOptionNodeDto<UnitDto> {
  Type: "Unit";
}

interface BaseOptionNodeDto<T extends OptionDto> {
  Type: string;
  Message: string;
  Data: DatumDto<T>[];
}

interface DatumDto<T extends OptionDto> {
  Option: T;
  Children: OptionNodeDto[];
}

export const optionNodeToModel = <T extends OptionNodeDto>(
  optionNode: T,
): OptionNode => {
  const optionConverter: any = {
    PieceType: pieceTypeToModel,
    Square: squareToModel,
    Move: moveToModel,
    Unit: (_: UnitDto): Unit => ({}),
  }[optionNode.Type];
  return {
    type: optionNode.Type,
    message: optionNode.Message,
    data: optionNode.Data.map((datum) => ({
      option: optionConverter(datum.Option),
      children: datum.Children.map(optionNodeToModel),
    })),
  };
};

export type RouteDto = RouteItemDto[];

type RouteItemDto = OptionDto & {
  Type: string;
};

export const routeToDto = (route: Route): RouteDto => route.map(routeItemToDto);

const routeItemToDto = <T extends OptionNode>({
  node,
  datum,
}: RouteItem<T>): RouteItemDto => {
  const optionConverter: any = {
    PieceType: pieceTypeToDto,
    Square: (square: Square) => ({ Square: squareToDto(square) }),
    Move: moveToDto,
    Unit: (_: Unit): UnitDto => ({}),
  }[node.type];
  return {
    Type: node.type,
    ...optionConverter(datum.option),
  };
};
