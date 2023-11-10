import { Move } from "@/model/game/move";
import {
  MoveOptionNode,
  NodeGroup,
  OptionDatum,
  OptionNode,
  Route,
} from "@/model/game/options";
import { Square } from "@/model/game/square";
import { ReactNode, createContext, useContext, useEffect, useState } from "react";

export const OptionContext = createContext<OptionContextValue>(null!);

export function useOptions() {
  return useContext(OptionContext);
}

export interface OptionContextValue {
  route: Route;
  isDone: boolean;
  moveMap: MoveMap;
  choose: (datum: OptionDatum) => void;
}

type MoveMap = {
  [from: string]: {
    [to: string]: NodeGroup<Move>;
  };
};

export const OptionProvider = ({
  root,
  children,
}: {
  root?: OptionNode;
  children?: ReactNode;
}) => {
  const isReady = root !== undefined;

  const [current, setCurrent] = useState<OptionNode[]>([]);
  const [route, setRoute] = useState<Route>([]);

  useEffect(() => {
    setCurrent(isReady ? [root] : [])
    setRoute([])
  }, [isReady, root, setCurrent, setRoute])

  const isDone = current.length == 0;
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
  const choose = (datum: OptionDatum) => {
    setCurrent(datum.children);
    setRoute([...route, datum.option]);
  };

  return (
    <OptionContext.Provider value={{ route, isDone, moveMap, choose }}>
      {children}
    </OptionContext.Provider>
  );
};

const is =
  <T extends OptionNode>(type: T["type"]) =>
  (node: OptionNode): node is T =>
    node.type == type;
