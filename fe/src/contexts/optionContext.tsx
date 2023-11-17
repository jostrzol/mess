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

export const useOptions = () => {
  return useContext(OptionContext);
};

export interface OptionContextValue {
  isReady: boolean;
  route: Route;
  currentNodes: OptionNode[];
  selectedNode: OptionNode | null;
  moveMap: MoveMap;
  choose: <T extends OptionNode>(node: T, datum: T["data"][number]) => void;
  select: <T extends OptionNode>(node: T) => void;
  reset: () => void;
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
  root: OptionNode | null;
  onChooseFinish?: (route: Route) => void;
  children?: ReactNode;
}) => {
  const isReady = root !== null;

  const [route, setRoute] = useState<Route>([]);
  const [current, setCurrent] = useState<OptionNode[]>([]);
  const [selected, setSelected] = useState<OptionNode | null>(null);

  useEffect(() => {
    setRoute([]);
    const newCurrent = isReady ? [root] : [];
    setCurrent(newCurrent);
    setSelected(newCurrent[0] ?? null);
  }, [isReady, root, setCurrent, setRoute]);

  const moveMap =
    selected?.type === "Move"
      ? selected.data.reduce((map, datum) => {
          const from = Square.toString(datum.option.from);
          const subMap = map[from] ?? {};
          const to = Square.toString(datum.option.to);
          const group = subMap[to] ?? [];
          const element = { node: selected, datum };
          return { ...map, [from]: { ...subMap, [to]: [...group, element] } };
        }, {} as MoveMap)
      : {};

  const choose = <T extends OptionNode>(node: T, datum: T["data"][number]) => {
    const newCurrent = datum.children;
    const newRouteItem: RouteItem<T> = [node, datum.option];
    const newRoute = [...route, newRouteItem];

    setRoute(newRoute);
    setCurrent(newCurrent);
    setSelected(newCurrent[0] ?? null);

    if (newCurrent.length === 0) {
      onChooseFinish?.(newRoute);
    }
  };

  const reset = () => {
    setRoute([]);
    const newCurrent = isReady ? [root] : [];
    setCurrent(newCurrent);
    setSelected(newCurrent[0] ?? null);
  };

  const select = <T extends OptionNode>(node: T) => {
    setSelected(node);
  };

  return (
    <OptionContext.Provider
      value={{
        isReady,
        route,
        currentNodes: current,
        selectedNode: selected,
        moveMap,
        choose,
        select,
        reset,
      }}
    >
      {children}
    </OptionContext.Provider>
  );
};
