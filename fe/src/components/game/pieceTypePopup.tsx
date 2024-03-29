import { useGameState } from "@/contexts/gameStateContext";
import { useOptions } from "@/contexts/optionContext";
import { useStaticData } from "@/contexts/staticDataContext";
import { Popup } from "../popup";
import { PieceIcon } from "./piece";

export const PieceTypePopup = () => {
  const { myColor } = useStaticData();
  const { isMyTurn } = useGameState();
  const { selectedNode, choose } = useOptions();
  if (!isMyTurn || selectedNode?.type !== "PieceType") {
    return null;
  }
  return (
    <Popup title={selectedNode.message} position="bottom">
      <div className="grid grid-flow-col auto-cols-fr gap-4 items-center">
        {selectedNode?.data.map((datum, i) => {
          const pieceType = datum.option;
          return (
            <div
              key={i}
              className="mx-auto flex flex-col align-center hover:scale-110 cursor-pointer w-max"
            >
              <div
                className="w-12 h-12"
                onClick={() => choose({ node: selectedNode, datum })}
              >
                <PieceIcon
                  blockRotation
                  color={myColor}
                  presentation={pieceType.presentation[myColor]}
                />
              </div>
              <p className="text-xs text-center select-none">
                {pieceType.name}
              </p>
            </div>
          );
        })}
      </div>
    </Popup>
  );
};
