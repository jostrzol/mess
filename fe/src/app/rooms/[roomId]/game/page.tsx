"use client";

import { GameApi } from "@/api/game";
import { GameChanged } from "@/api/schema/event";
import { ConnectionStatus } from "@/components/connectionStatus";
import { Board } from "@/components/game/board";
import { OptionIndicator } from "@/components/game/optionIndicator";
import { ResolutionPopup } from "@/components/game/resolutionPopup";
import { GameStateProvider } from "@/contexts/gameStateContext";
import { useMessApi } from "@/contexts/messApiContext";
import { OptionProvider } from "@/contexts/optionContext";
import { useRoomWebsocket } from "@/contexts/roomWsContext";
import { StaticDataProvider } from "@/contexts/staticDataContext";
import { Route } from "@/model/game/options";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { RoomPageParams } from "../layout";
import {PieceTypePopup} from "@/components/game/pieceTypePopup";
import {UnitPopup} from "@/components/game/unitPopup";

const GamePage = ({ params }: RoomPageParams) => {
  const gameApi = useMessApi(GameApi);
  const client = useQueryClient();

  const keyGame = ["room", params.roomId, "game"];

  const keyStaticData = [...keyGame, "static"];
  const { data: staticData } = useQuery({
    queryKey: keyStaticData,
    queryFn: () => gameApi.getStaticData(params.roomId),
    staleTime: Infinity,
  });

  const keyDynamic = [...keyGame, "dynamic"];

  const keyState = [...keyDynamic, "state"];
  const { data: state } = useQuery({
    queryKey: keyState,
    queryFn: () => gameApi.getState(params.roomId),
  });

  const keyOptions = [...keyDynamic, "options"];
  const { data: optionTree } = useQuery({
    queryKey: keyOptions,
    queryFn: () => gameApi.getTurnOptions(params.roomId),
  });

  const keyResolution = [...keyDynamic, "resolution"];
  const {
    data: resolution,
  } = useQuery({
    queryKey: keyResolution,
    queryFn: () => gameApi.getResolution(params.roomId),
  });
  const status = resolution?.status ?? "Unresolved";

  const { mutate } = useMutation({
    mutationKey: ["room", params.roomId, "game", "turn"],
    mutationFn: (route: Route) =>
      gameApi.playTurn(params.roomId, state!.turnNumber, route),
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
    <StaticDataProvider roomId={params.roomId} staticData={staticData}>
      <GameStateProvider state={state}>
        <OptionProvider
          root={optionTree ?? null}
          onChooseFinish={(route) => mutate(route)}
        >
          <ConnectionStatus />
          <div className="pt-4">
            <OptionIndicator />
          </div>
          <Board board={staticData.board} />
          <PieceTypePopup />
          <UnitPopup />
          {status !== "Unresolved" && <ResolutionPopup status={status} />}
        </OptionProvider>
      </GameStateProvider>
    </StaticDataProvider>
  );
};

export default GamePage;
