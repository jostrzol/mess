import { Room } from "@/model/room";
import { UUID } from "crypto";
import useWebSocket from "react-use-websocket";
import { url } from "./url";
import { throwIfError } from "./utils";

interface RoomDto {
  ID: UUID;
  Players: number;
}

const toModel = (room: RoomDto): Room => {
  return new Room(room.ID, room.Players);
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

interface RoomChanged {
  EventType: "RoomChanged";
}

interface Moved {
  EventType: "Moved";
}

type Event = RoomChanged | Moved;

export const useRoomWebsocket = (roomId: UUID) => {
  const url_ = url("rooms/:id/websocket", {
    params: { id: roomId },
    schema: "ws",
  });
  const { lastJsonMessage, ...rest } = useWebSocket(url_);
  const lastEvent = lastJsonMessage as Event | null;
  return { lastEvent, ...rest };
};
