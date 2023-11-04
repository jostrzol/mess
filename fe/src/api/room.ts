import { Room } from "@/model/room";
import { UUID } from "crypto";
import useWebSocket from "react-use-websocket";
import { url } from "./url";
import { throwIfError } from "./utils";

interface RoomDto {
  ID: UUID;
  Players: number;
  PlayersNeeded: number;
  IsStartable: boolean;
  IsStarted: boolean;
}

const toModel = (room: RoomDto): Room => {
  return {
    id: room.ID,
    players: room.Players,
    playersNeeded: room.PlayersNeeded,
    isStartable: room.IsStartable,
    isStarted: room.IsStarted,
  };
};

export const createRoom = async (): Promise<Room> => {
  const res = await fetch(url("rooms"), {
    method: "POST",
    credentials: "include",
  });
  await throwIfError(res);

  const obj: RoomDto = await res.json();
  return toModel(obj);
};

export const joinRoom = async (id: UUID): Promise<Room> => {
  const url_ = url("rooms/:id/players", { params: { id: id } });
  const res = await fetch(url_, { method: "PUT", credentials: "include" });
  await throwIfError(res);

  const obj: RoomDto = await res.json();
  return toModel(obj);
};

export const startGame = async (id: UUID): Promise<Room> => {
  const url_ = url("rooms/:id/game", { params: { id: id } });
  const res = await fetch(url_, { method: "PUT", credentials: "include" });
  await throwIfError(res);

  const obj: RoomDto = await res.json();
  return toModel(obj);
};

interface RoomChanged {
  EventType: "RoomChanged";
}

interface GameStarted {
  EventType: "GameStarted";
}

type Event = RoomChanged | GameStarted;

export const useRoomWebsocket = (roomId: UUID) => {
  const url_ = url("rooms/:id/websocket", {
    params: { id: roomId },
    schema: "ws",
  });
  const { lastJsonMessage, ...rest } = useWebSocket(url_);
  const lastEvent = lastJsonMessage as Event | null;
  return { lastEvent, ...rest };
};
