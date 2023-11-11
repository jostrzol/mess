"use client";

import { getState, getStaticData, getTurnOptions, playTurn } from "@/api/game";
import { GameChanged } from "@/api/schema/event";
import { Board } from "@/components/game/board";
import { GameStateProvider } from "@/contexts/gameStateContext";
import { OptionProvider } from "@/contexts/optionContext";
import { useRoomWebsocket } from "@/contexts/roomWsContext";
import { StaticDataProvider } from "@/contexts/staticDataContext";
import { Route } from "@/model/game/options";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { RoomPageParams } from "../layout";

const GamePage = ({ params }: RoomPageParams) => {
  const client = useQueryClient();

  const keyStaticData = ["room", params.roomId, "game", "static"];
  const { data: staticData } = useQuery({
    queryKey: keyStaticData,
    queryFn: () => getStaticData(params.roomId),
    staleTime: Infinity,
  });

  const keyState = ["room", params.roomId, "game", "state"];
  const { data: state } = useQuery({
    queryKey: keyState,
    queryFn: () => getState(params.roomId),
  });

  const keyOptions = ["room", params.roomId, "game", "options"];
  const { data: optionTree } = useQuery({
    queryKey: keyOptions,
    queryFn: () => getTurnOptions(params.roomId),
  });

  const { mutate } = useMutation({
    mutationKey: ["room", params.roomId, "game", "turn"],
    mutationFn: (route: Route) =>
      playTurn(params.roomId, state!.turnNumber, route),
    onSuccess: (newState) => {
      client.setQueryData(keyState, newState);
      client.invalidateQueries({ queryKey: keyOptions });
    },
    onError: (e) => {
      console.error(e);
      client.invalidateQueries({ queryKey: ["room", params.roomId, "game"] });
    },
  });
  useRoomWebsocket<GameChanged>({
    type: "GameChanged",
    onEvent: (e) => {
      console.log(e);
      client.invalidateQueries({ queryKey: ["room", params.roomId, "game"] });
      client.resetQueries({ queryKey: keyOptions });
    },
  });

  if (staticData === undefined || state === undefined) {
    return null;
  }
  return (
    <StaticDataProvider staticData={staticData}>
      <GameStateProvider state={state}>
        <OptionProvider
          root={optionTree ?? null}
          onChooseFinish={(route) => {
            mutate(route);
          }}
        >
          <Board board={staticData.board} />
        </OptionProvider>
      </GameStateProvider>
    </StaticDataProvider>
  );
};

export default GamePage;
