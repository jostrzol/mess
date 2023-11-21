import { useGameState } from "@/contexts/gameStateContext";
import { useOptions } from "@/contexts/optionContext";
import { MdPlayArrow } from "react-icons/md";
import { Popup } from "../popup";
import clsx from "clsx";

export const UnitPopup = () => {
  const { isMyTurn } = useGameState();
  const { selectedNode, choose } = useOptions();
  if (!isMyTurn || selectedNode?.type !== "Unit") {
    return null;
  }
  const datum = selectedNode.data[0];
  return (
    <Popup title={selectedNode.message} position="bottom">
      <div
        className={clsx(
          "cursor-pointer text-success-strong mx-auto",
          "transition-transform hover:translate-x-1",
        )}
        onClick={() => choose({ node: selectedNode, datum })}
      >
        <MdPlayArrow size={64} />
      </div>
    </Popup>
  );
};
