"use client";

import { getGame, playTurn } from "@/api/game";
import { Board } from "@/components/game/board";
import { GameStateProvider } from "@/contexts/gameStateContext";
import { OptionProvider } from "@/contexts/optionContext";
import { Route } from "@/model/game/options";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { RoomPageParams } from "../layout";

const GamePage = ({ params }: RoomPageParams) => {
  const client = useQueryClient();
  const { data: state, isSuccess } = useQuery({
    queryKey: ["room", params.roomId, "game"],
    queryFn: () => getGame(params.roomId),
  });
  const { mutate } = useMutation({
    mutationKey: ["room", params.roomId, "game", "turn"],
    mutationFn: (route: Route) =>
      playTurn(params.roomId, state!.turnNumber, route),
    onSuccess: (newState) => {
      client.setQueryData(["room", params.roomId, "game"], newState);
    },
    onError: (e) => {
      console.error(e)
      client.invalidateQueries({queryKey: ["room", params.roomId, "game"]})
    },
  });
  if (!isSuccess) {
    return null;
  }
  return (
    <GameStateProvider state={state}>
      <OptionProvider
        root={state.optionTree}
        onChooseFinish={(route) => {
          mutate(route);
        }}
      >
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
