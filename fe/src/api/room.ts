import { Room } from "@/model/room";
import { UUID } from "crypto";
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
    cache: "no-store",
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
