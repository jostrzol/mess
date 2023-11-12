"use client";

import { getResolution, getState, getStaticData, getTurnOptions, playTurn } from "@/api/game";
import { GameChanged } from "@/api/schema/event";
import { Board } from "@/components/game/board";
import { GameStateProvider } from "@/contexts/gameStateContext";
import { OptionProvider } from "@/contexts/optionContext";
import { useRoomWebsocket } from "@/contexts/roomWsContext";
import { StaticDataProvider } from "@/contexts/staticDataContext";
import { Route } from "@/model/game/options";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { RoomPageParams } from "../layout";
import {Resolution} from "@/model/game/resolution";
import {ResolutionPopup} from "@/components/game/resolutionPopup";
import {Options} from "@/components/game/options";

const GamePage = ({ params }: RoomPageParams) => {
  const client = useQueryClient();

  const keyGame = ["room", params.roomId, "game"]

  const keyStaticData = [...keyGame, "static"];
  const { data: staticData } = useQuery({
    queryKey: keyStaticData,
    queryFn: () => getStaticData(params.roomId),
    staleTime: Infinity,
  });

  const keyDynamic = [...keyGame, "dynamic"];

  const keyState = [...keyDynamic, "state"];
  const { data: state } = useQuery({
    queryKey: keyState,
    queryFn: () => getState(params.roomId),
  });

  const keyOptions = [...keyDynamic, "options"];
  const { data: optionTree } = useQuery({
    queryKey: keyOptions,
    queryFn: () => getTurnOptions(params.roomId),
  });

  const keyResolution = [...keyDynamic, "resolution"];
  const { data: { status } } = useQuery({
    queryKey: keyResolution,
    queryFn: () => getResolution(params.roomId),
    initialData: { status: "Unresolved" } as Resolution
  });

  const { mutate } = useMutation({
    mutationKey: ["room", params.roomId, "game", "turn"],
    mutationFn: (route: Route) =>
      playTurn(params.roomId, state!.turnNumber, route),
    onSuccess: (newState) => {
      client.invalidateQueries({ queryKey: keyDynamic });
      client.setQueryData(keyState, newState);
    },
    onError: () => {
      client.invalidateQueries({ queryKey: keyDynamic });
    },
  });
  useRoomWebsocket<GameChanged>({
    type: "GameChanged",
    onEvent: () => {
      client.invalidateQueries({ queryKey: keyDynamic });
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
          <Options />
          <Board board={staticData.board} />
          {status !== "Unresolved" && <ResolutionPopup status={status}/>}
        </OptionProvider>
      </GameStateProvider>
    </StaticDataProvider>
  );
};

export default GamePage;
