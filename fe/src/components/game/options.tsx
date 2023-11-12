import { useOptions } from "@/contexts/optionContext";
import { OptionNode } from "@/model/game/options";
import clsx from "clsx";

export const Options = () => {
  const { currentNodes, selectedNode, select } = useOptions();
  return (
    <div className="relative h-10">
      <div className="grid grid-flow-col auto-cols-fr">
        {[...currentNodes].map((node, i) => (
          <Option
            key={i}
            optionNode={node}
            isSelected={node === selectedNode}
            onClick={select}
          />
        ))}
      </div>
    </div>
  );
};

const Option = ({
  optionNode,
  isSelected,
  onClick,
}: {
  optionNode: OptionNode;
  isSelected: boolean;
  onClick?: (node: OptionNode) => void;
}) => {
  return (
    <div
      className={clsx(
        "mx-2 w-full",
        "flex flex-col align-center",
        "[&:first-child_hr.left]:invisible",
        "[&:last-child_hr.right]:invisible",
      )}
    >
      <div className="relative w-full">
        <div
          className={clsx(
            "w-5 h-5 mx-auto border-2 border-primary rounded-full bg-background group",
            onClick && "cursor-pointer",
          )}
          onClick={() => onClick?.(optionNode)}
        >
          <div
            className={clsx(
              "w-full h-full p-[0.4rem] rounded-full",
              "bg-clip-content bg-primary",
              "transition-transform",
              isSelected ? "scale-[275%]" : "group-hover:scale-[200%]"
            )}
          />
        </div>
        <hr className="left absolute w-1/2 top-1/2 left-0 -z-10 -translate-y-1/2 border-b-2 border-primary" />
        <hr className="right absolute w-1/2 top-1/2 right-0 -z-10 -translate-y-1/2 border-b-2 border-primary" />
        <div className="absolute w-full h-full bg-background -z-10" />
      </div>
      <p className="text-sm text-center select-none">{optionNode.message}</p>
    </div>
  );
};
