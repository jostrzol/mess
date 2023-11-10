import {
  MoveOptionNode,
  NodeGroup,
  OptionNode,
  Route,
  RouteItem,
} from "@/model/game/options";
import { Square } from "@/model/game/square";
import {
  ReactNode,
  createContext,
  useContext,
  useEffect,
  useState,
} from "react";

export const OptionContext = createContext<OptionContextValue>(null!);

export function useOptions() {
  return useContext(OptionContext);
}

export interface OptionContextValue {
  route: Route;
  isDone: boolean;
  moveMap: MoveMap;
  choose: <T extends OptionNode>(node: T, datum: T["data"][number]) => void;
}

type MoveMap = {
  [from: string]: {
    [to: string]: NodeGroup<MoveOptionNode>;
  };
};

export const OptionProvider = ({
  root,
  onChooseFinish,
  children,
}: {
  root?: OptionNode;
  onChooseFinish?: (route: Route) => void;
  children?: ReactNode;
}) => {
  const isReady = root !== undefined;

  const [current, setCurrent] = useState<OptionNode[]>([]);
  const [route, setRoute] = useState<Route>([]);

  useEffect(() => {
    setCurrent(isReady ? [root] : []);
    setRoute([]);
  }, [isReady, root, setCurrent, setRoute]);

  const isDone = current.length == 0;
  const moveMap = current
    .filter(is<MoveOptionNode>("Move"))
    .reduce((map, node) => {
      return node.data.reduce((map, datum) => {
        const from = Square.toString(datum.option.from);
        const subMap = map[from] ?? {};
        const to = Square.toString(datum.option.to);
        const group = subMap[to] ?? [];
        const element = { node: node, datum };
        return { ...map, [from]: { ...subMap, [to]: [...group, element] } };
      }, map);
    }, {} as MoveMap);
  const choose = <T extends OptionNode>(node: T, datum: T["data"][number]) => {
    const newCurrent = datum.children;
    const newRouteItem: RouteItem<T> = [node, datum.option];
    const newRoute = [...route, newRouteItem];

    setCurrent(newCurrent);
    setRoute(newRoute);

    if (newCurrent.length === 0) {
      onChooseFinish?.(newRoute);
    }
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
