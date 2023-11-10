"use client";

import { getGame } from "@/api/game";
import { Board } from "@/components/game/board";
import { GameStateProvider } from "@/contexts/gameStateContext";
import { OptionProvider } from "@/contexts/optionContext";
import { useQuery } from "@tanstack/react-query";
import { RoomPageParams } from "../layout";

const GamePage = ({ params }: RoomPageParams) => {
  const { data: state, isSuccess } = useQuery({
    queryKey: ["room", params.roomId, "game"],
    queryFn: () => getGame(params.roomId),
  });
  if (!isSuccess) {
    return null;
  }
  return (
    <GameStateProvider state={state}>
      <OptionProvider root={state.optionTree}>
        <Board
          board={{
            height: 8,
            width: 8,
          }}
        />
      </OptionProvider>
    </GameStateProvider>
  );
};

export default GamePage;
