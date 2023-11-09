import { PieceType } from "@/model/game/pieceType";
import { Square } from "@/model/game/square";
import { useEffect, useState } from "react";
import { Move } from "./move";

export type OptionNode =
  | PieceTypeOptionNode
  | SquareOptionNode
  | MoveOptionNode
  | UnitOptionNode;

export type Route = Option[];

export type Option = PieceType | Square | Move | Unit;

export type Unit = [];

export interface PieceTypeOptionNode extends BaseOptionNode<PieceType> {
  type: "PieceType";
}

export interface SquareOptionNode extends BaseOptionNode<Square> {
  type: "Square";
}

export interface MoveOptionNode extends BaseOptionNode<Move> {
  type: "Move";
}

export interface UnitOptionNode extends BaseOptionNode<Unit> {
  type: "Unit";
}

interface BaseOptionNode<T extends Option> {
  type: string;
  message: string;
  data: Datum<T>[];
}

interface Datum<T extends Option> {
  option: T;
  children: OptionNode[];
}

const is =
  <T extends OptionNode>(type: T["type"]) =>
  (node: OptionNode): node is T =>
    node.type == type;

type NodeGroup<T extends Option> = { message: string; datum: Datum<T> }[];

export type MoveMap = {
  [from: string]: {
    [to: string]: NodeGroup<Move>;
  };
};

export interface UseOptionTreeValue {
  route: Route;
  current: OptionNode[];
  choose: (datum: Datum<Option>) => void;
  isDone: boolean;
  moveMap: MoveMap;
}

export const useOptionTree = (root?: OptionNode): UseOptionTreeValue => {
  const isInit = root !== undefined;

  const [route, setRoute] = useState<Route>([]);
  const [current, setCurrent] = useState(isInit ? [root] : []);

  useEffect(() => {
    setCurrent(isInit ? [root] : []);
  }, [isInit, root]);

  const choose = (datum: Datum<Option>) => {
    if (!isInit) {
      throw new Error("tried to choose option on uninitialized option tree");
    }
    setCurrent(datum.children);
    setRoute([...route, datum.option]);
  };
  const isDone = isInit && current.length == 0;

  const moveMap = current
    .filter(is<MoveOptionNode>("Move"))
    .reduce((map, node) => {
      return node.data.reduce((map, datum) => {
        const from = Square.toString(datum.option.from);
        const subMap = map[from] ?? {};
        const to = Square.toString(datum.option.to);
        const group = subMap[to] ?? [];
        const element = { message: node.message, datum };
        return { ...map, [from]: { ...subMap, [to]: [...group, element] } };
      }, map);
    }, {} as MoveMap);

  return { route, current, choose, isDone, moveMap };
};
