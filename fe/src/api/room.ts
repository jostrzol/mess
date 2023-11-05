import { Room } from "@/model/room";
import { UUID } from "crypto";
import useWebSocket from "react-use-websocket";
import { url } from "./url";
import { throwIfError } from "./utils";
import {RoomDto, roomToModel} from "./schema/room";
import {Event} from "./schema/event";

export const createRoom = async (): Promise<Room> => {
  const res = await fetch(url("rooms"), {
    method: "POST",
    credentials: "include",
  });
  await throwIfError(res);

  const obj: RoomDto = await res.json();
  return roomToModel(obj);
};

export const joinRoom = async (id: UUID): Promise<Room> => {
  const url_ = url("rooms/:id/players", { params: { id: id } });
  const res = await fetch(url_, { method: "PUT", credentials: "include" });
  await throwIfError(res);

  const obj: RoomDto = await res.json();
  return roomToModel(obj);
};

export const useRoomWebsocket = (roomId: UUID) => {
  const url_ = url("rooms/:id/websocket", {
    params: { id: roomId },
    schema: "ws",
  });
  const { lastJsonMessage, ...rest } = useWebSocket(url_);
  const lastEvent = lastJsonMessage as Event | null;
  return { lastEvent, ...rest };
};
