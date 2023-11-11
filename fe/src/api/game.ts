import { GameState } from "@/model/game/gameState";
import { StaticData } from "@/model/game/gameStaticData";
import { OptionNode, Route } from "@/model/game/options";
import { UUID } from "crypto";
import { GameStateDto, gameStateToModel } from "./schema/game";
import { OptionNodeDto, optionNodeToModel, routeToDto } from "./schema/options";
import { StaticDataDto, staticDataToModel } from "./schema/staticData";
import { url } from "./url";
import { throwIfError } from "./utils";

export const getStaticData = async (roomId: UUID): Promise<StaticData> => {
  const url_ = url("rooms/:id/game/static", { params: { id: roomId } });
  const res = await fetch(url_, { method: "GET", credentials: "include" });
  await throwIfError(res);

  const obj: StaticDataDto = await res.json();
  return staticDataToModel(obj);
};

export const getState = async (roomId: UUID): Promise<GameState> => {
  const url_ = url("rooms/:id/game/state", { params: { id: roomId } });
  const res = await fetch(url_, { method: "GET", credentials: "include" });
  await throwIfError(res);

  const obj: GameStateDto = await res.json();
  return gameStateToModel(obj);
};

export const getTurnOptions = async (roomId: UUID): Promise<OptionNode> => {
  const url_ = url("rooms/:id/game/options", { params: { id: roomId } });
  const res = await fetch(url_, { method: "GET", credentials: "include" });
  await throwIfError(res);

  const obj: OptionNodeDto = await res.json();
  return optionNodeToModel(obj);
};

export const playTurn = async (
  roomId: UUID,
  turn: number,
  route: Route,
): Promise<GameState> => {
  const url_ = url("rooms/:id/game/turns/:turn", {
    params: { id: roomId, turn: turn },
  });
  const res = await fetch(url_, {
    method: "PUT",
    credentials: "include",
    body: JSON.stringify(routeToDto(route)),
  });
  await throwIfError(res);

  const obj: GameStateDto = await res.json();
  return gameStateToModel(obj);
};
