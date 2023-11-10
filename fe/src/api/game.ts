import { GameState } from "@/model/game/gameState";
import { Room } from "@/model/room";
import { UUID } from "crypto";
import { GameStateDto, gameStateToModel } from "./schema/game";
import { RouteDto, routeToDto } from "./schema/options";
import { RoomDto, roomToModel } from "./schema/room";
import { url } from "./url";
import { throwIfError } from "./utils";
import {Route} from "@/model/game/options";

export const startGame = async (roomId: UUID): Promise<Room> => {
  const url_ = url("rooms/:id/game", { params: { id: roomId } });
  const res = await fetch(url_, { method: "PUT", credentials: "include" });
  await throwIfError(res);

  const obj: RoomDto = await res.json();
  return roomToModel(obj);
};

export const getGame = async (roomId: UUID): Promise<GameState> => {
  const url_ = url("rooms/:id/game", { params: { id: roomId } });
  const res = await fetch(url_, { method: "GET", credentials: "include" });
  await throwIfError(res);

  const obj: GameStateDto = await res.json();
  return gameStateToModel(obj);
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
