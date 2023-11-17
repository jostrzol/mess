import { PieceType } from "@/model/game/pieceType";
import { Square } from "@/model/game/square";
import { Move } from "./move";

export type OptionNode =
  | PieceTypeOptionNode
  | SquareOptionNode
  | MoveOptionNode
  | UnitOptionNode;

export type OptionDatum = OptionNode["data"][number];

export type RouteItem<T extends OptionNode> = [T, T["data"][number]["option"]];

export type Route = RouteItem<OptionNode>[];

export type Option = PieceType | Square | Move | Unit;

export type Unit = {};

export type PieceTypeOptionNode = BaseOptionNode<"PieceType", PieceType>;

export type SquareOptionNode = BaseOptionNode<"Square", Square>;

export type MoveOptionNode = BaseOptionNode<"Move", Move>;

export type UnitOptionNode = BaseOptionNode<"Unit", Unit>;

interface BaseOptionNode<T extends string, TD extends Option> {
  type: T;
  message: string;
  data: BaseDatum<TD>[];
}

interface BaseDatum<T extends Option> {
  option: T;
  children: OptionNode[];
}

export type NodeGroup<T extends OptionNode> = {
  node: BaseOptionNode<T["type"], T["data"][number]["option"]>;
  datum: T["data"][number];
}[];
