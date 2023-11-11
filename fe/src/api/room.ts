import { Room } from "@/model/room";
import { UUID } from "crypto";
import { RoomDto, roomToModel } from "./schema/room";
import { url } from "./url";
import { throwIfError } from "./utils";

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

export const startGame = async (roomId: UUID): Promise<Room> => {
  const url_ = url("rooms/:id/game", { params: { id: roomId } });
  const res = await fetch(url_, { method: "PUT", credentials: "include" });
  await throwIfError(res);

  const obj: RoomDto = await res.json();
  return roomToModel(obj);
};
