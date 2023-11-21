import {
  MoveOptionNode,
  OptionNode,
  Route,
  RouteItem,
  SquareOptionNode,
  UnitOptionNode,
} from "@/model/game/options";
import { Square } from "@/model/game/square";
import {
  ReactNode,
  createContext,
  useCallback,
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
  squareMap: SquareMap;
  choose: <T extends OptionNode>(routeItem: RouteItem<T>) => void;
  select: <T extends OptionNode>(node: T) => void;
  isResetable: boolean;
  reset: () => void;
}

type MoveMap = {
  [from: string]: {
    [to: string]: RouteItem<MoveOptionNode>;
  };
};

type SquareMap = { [square: string]: RouteItem<SquareOptionNode> };

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
  const [isResetable, setIsResetable] = useState<boolean>(false);

  const reset = useCallback(() => {
    setRoute([]);
    const newCurrent = isReady ? [root] : [];
    setCurrent(newCurrent);
    setSelected(newCurrent[0] ?? null);
    setIsResetable(false);
  }, [isReady, root]);

  useEffect(reset, [reset]);

  const choose = useCallback(
    <T extends OptionNode>(routeItem: RouteItem<T>) => {
      const newCurrent = routeItem.datum.children;
      const newRoute = [...route, routeItem];

      setRoute(newRoute);
      setCurrent(newCurrent);
      setSelected(newCurrent[0] ?? null);
      setIsResetable(true);

      if (newCurrent.length === 0) {
        onChooseFinish?.(newRoute);
      }
    },
    [onChooseFinish, route],
  );

  useEffect(() => {
    const routeItem = singleUnitWithSingleChild(current);
    if (routeItem) {
      const lastCanReset = isResetable;
      choose(routeItem);
      setIsResetable(lastCanReset || false);
    }
  }, [current, choose, isResetable]);

  const select = <T extends OptionNode>(node: T) => setSelected(node);

  const moveMap =
    selected?.type === "Move"
      ? selected.data.reduce((map, datum) => {
          const from = Square.toString(datum.option.from);
          const subMap = map[from] ?? {};
          const to = Square.toString(datum.option.to);
          const routeItem = { node: selected, datum };
          return { ...map, [from]: { ...subMap, [to]: routeItem } };
        }, {} as MoveMap)
      : {};

  const squareMap =
    selected?.type === "Square"
      ? selected.data.reduce((map, datum) => {
          const from = Square.toString(datum.option);
          const routeItem = {node: selected, datum}
          return { ...map, [from]: routeItem};
        }, {} as SquareMap)
      : {};

  return (
    <OptionContext.Provider
      value={{
        isReady,
        route,
        currentNodes: current,
        selectedNode: selected,
        moveMap,
        squareMap,
        choose,
        select,
        isResetable,
        reset,
      }}
    >
      {children}
    </OptionContext.Provider>
  );
};

const singleUnitWithSingleChild = (
  nodes: OptionNode[],
): RouteItem<UnitOptionNode> | null => {
  if (nodes.length != 1) return null;

  const node = nodes[0];
  if (node.type == "Unit" && node.data.length == 1) {
    return { node, datum: node.data[0] };
  }

  return null;
};
