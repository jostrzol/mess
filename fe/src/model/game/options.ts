import { PieceType } from "@/model/game/pieceType";
import { Square } from "@/model/game/square";
import { Dispatch, useReducer } from "react";
import { Move } from "./move";

export type OptionNode =
  | PieceTypeOptionNode
  | SquareOptionNode
  | MoveOptionNode
  | UnitOptionNode;

export type OptionDatum = OptionNode["data"][number];

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
  data: BaseDatum<T>[];
}

interface BaseDatum<T extends Option> {
  option: T;
  children: OptionNode[];
}

export type NodeGroup<T extends Option> = { message: string; datum: BaseDatum<T> }[];

// const is =
//   <T extends OptionNode>(type: T["type"]) =>
//   (node: OptionNode): node is T =>
//     node.type == type;
//
//
// export type MoveMap = {
//   [from: string]: {
//     [to: string]: NodeGroup<Move>;
//   };
// };
//
// export interface UseOptionTreeValue {
//   route: Route;
//   isDone: boolean;
//   moveMap: MoveMap;
//   dispatch: Dispatch<Action>;
// }
//
// interface SetRoot {
//   type: "SetRoot";
//   root?: OptionNode;
// }
//
// interface Choose {
//   type: "Choose";
//   datum: OptionDatum;
// }
//
// export type Action = SetRoot | Choose;
//
// interface State {
//   root?: OptionNode;
//   current: OptionNode[];
//   route: Option[];
// }
//
// export const useOptionTree = (): UseOptionTreeValue => {
//   const [state, dispatch] = useReducer(
//     (state: State, action: Action) => {
//       switch (action.type) {
//         case "SetRoot":
//           return {
//             root: action.root,
//             current: action.root === undefined ? [] : [action.root],
//             route: [],
//           };
//         case "Choose":
//           return {
//             ...state,
//             current: action.datum.children,
//             route: [...state.route, action.datum.option],
//           };
//         default:
//           return state;
//       }
//     },
//     {
//       root: undefined,
//       current: [],
//       route: [],
//     },
//   );
//
//   const { current } = state;
//
//   const isDone = current.length == 0;
//   const moveMap = current
//     .filter(is<MoveOptionNode>("Move"))
//     .reduce((map, node) => {
//       return node.data.reduce((map, datum) => {
//         const from = Square.toString(datum.option.from);
//         const subMap = map[from] ?? {};
//         const to = Square.toString(datum.option.to);
//         const group = subMap[to] ?? [];
//         const element = { message: node.message, datum };
//         return { ...map, [from]: { ...subMap, [to]: [...group, element] } };
//       }, map);
//     }, {} as MoveMap);
//
//   return { ...state, isDone, moveMap, dispatch };
// };
