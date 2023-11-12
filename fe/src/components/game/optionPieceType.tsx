import { Window } from "@/components/window";
import { useGameState } from "@/contexts/gameStateContext";
import { useOptions } from "@/contexts/optionContext";
import Image from "next/image";

export const OptionPieceType = () => {
  const { isMyTurn } = useGameState();
  const { selectedNode, choose } = useOptions();
  if (!isMyTurn || selectedNode?.type !== "PieceType") {
    return null;
  }
  return (
    <Window
      title={selectedNode.message}
      opaque
      className="fixed bottom-0 m-4 max-w-[90%] "
    >
      <div className="grid grid-flow-col auto-cols-fr gap-4">
        {selectedNode?.data.map((datum, i) => {
          const pieceType = datum.option;
          return (
            <div
              key={i}
              className="flex flex-col align-center hover:scale-110 cursor-pointer"
            >
              <Image
                className="mx-auto"
                alt={pieceType.name}
                src={pieceType.iconUri}
                height={64}
                width={64}
                onClick={() => choose(selectedNode, datum)}
              />
              <p className="text-xs text-center select-none">
                {pieceType.name}
              </p>
            </div>
          );
        })}
      </div>
    </Window>
  );
};
