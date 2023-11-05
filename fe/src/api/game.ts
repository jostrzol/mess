import {Room} from "@/model/room";
import {UUID} from "crypto";
import {url} from "./url";
import {throwIfError} from "./utils";
import {RoomDto, roomToModel} from "./schema/room";
import {GameState} from "@/model/gameState";
import {GameStateDto, gameStateToModel} from "./schema/game";

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
