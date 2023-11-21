import { useGameState } from "@/contexts/gameStateContext";
import { useOptions } from "@/contexts/optionContext";
import { useTheme } from "@/contexts/themeContext";
import { OptionNode } from "@/model/game/options";
import clsx from "clsx";
import { ReactNode } from "react";
import { MdClose, MdHourglassEmpty, MdRefresh } from "react-icons/md";

export const OptionIndicator = () => {
  const {
    theme: { colors },
  } = useTheme();
  const { isMyTurn } = useGameState();
  const { isReady, route, currentNodes, selectedNode, select, isResetable, reset } =
    useOptions();

  if (!isMyTurn) {
    return (
      <Options>
        <Indicator label={"Opponent's turn"}>
          <div className="w-fit h-fit text-primary hover:animate-spin-slow">
            <MdHourglassEmpty />
          </div>
        </Indicator>
      </Options>
    );
  } else if (!isReady) {
    return (
      <Options>
        <Indicator label={"Loading..."}>
          <div className="w-fit h-fit mt-[1px] ml-[1px] text-primary animate-spin-slow">
            <MdRefresh />
          </div>
        </Indicator>
      </Options>
    );
  }

  return (
    <Options>
      {isResetable && (
        <Indicator
          key="cancel"
          label={"Cancel"}
          circleColor={colors["danger-dim"]}
          onClick={reset}
        >
          <div className="w-fit h-fit ml-[1px] antialiased hover:animate-spin-1/4 text-danger-dim">
            <MdClose />
          </div>
        </Indicator>
      )}
      {currentNodes
        .map((node, i) => (
          <Option
            key={i}
            optionNode={node}
            selected={node === selectedNode}
            onClick={select}
          />
        ))
        .concat(route.length === 0 ? [] : [])}
    </Options>
  );
};

const Options = ({ children }: { children: ReactNode }) => (
  <div className="relative h-10">
    <div className="grid grid-flow-col auto-cols-fr">{children}</div>
  </div>
);

const Option = ({
  optionNode,
  selected,
  onClick,
}: {
  optionNode: OptionNode;
  selected: boolean;
  onClick?: (node: OptionNode) => void;
}) => {
  return (
    <Indicator
      label={optionNode.message}
      onClick={selected ? undefined : () => onClick?.(optionNode)}
    >
      <Dot selected={selected} />
    </Indicator>
  );
};

const Indicator = ({
  label,
  onClick,
  circleColor,
  circleBackground,
  className,
  children,
}: {
  label?: string;
  onClick?: () => void;
  circleColor?: string;
  circleBackground?: string;
  className?: string;
  children?: ReactNode;
}) => {
  return (
    <div
      className={clsx(
        "mx-2 w-full flex flex-col align-center",
        "[&:first-child_hr.left]:invisible",
        "[&:last-child_hr.right]:invisible",
        className,
      )}
    >
      <div className="relative w-full">
        <div
          className={clsx(
            "mx-auto w-5 h-5 border-2 border-primary rounded-full bg-background group",
            "flex items-center justify-center",
            onClick && "cursor-pointer",
          )}
          style={{
            borderColor: circleColor,
            backgroundColor: circleBackground,
          }}
          onClick={() => onClick?.()}
        >
          {children}
        </div>
        <Bar side="left" />
        <Bar side="right" />
        <div className="absolute w-full h-full bg-background -z-10" />
      </div>
      <p className="text-sm text-center select-none">{label}</p>
    </div>
  );
};

const Bar = ({ side }: { side: "left" | "right" }) => (
  <hr
    className={clsx(
      side,
      { left: "left-0", right: "right-0" }[side],
      "absolute w-1/2 top-1/2 -z-10 -translate-y-1/2 border-b-2 border-primary",
    )}
  />
);

const Dot = ({ selected = false }: { selected?: boolean }) => {
  return (
    <div
      className={clsx(
        "w-full h-full p-[0.4rem] rounded-full",
        "bg-clip-content bg-primary",
        "transition-transform",
        selected ? "scale-[275%]" : "group-hover:scale-[200%]",
      )}
    />
  );
};
