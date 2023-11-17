import { Room } from "@/model/room";
import { UUID } from "crypto";
import { MessApi } from "./messApi";
import { RoomDto, roomToModel } from "./schema/room";
import { throwIfError } from "./utils";

export class RoomApi extends MessApi {
  public createRoom = async (): Promise<Room> => {
    const res = await fetch(this.url("rooms"), {
      method: "POST",
      credentials: "include",
    });
    await throwIfError(res);

    const obj: RoomDto = await res.json();
    return roomToModel(obj);
  };

  public joinRoom = async (id: UUID): Promise<Room> => {
    const url_ = this.url("rooms/:id/players", { params: { id: id } });
    const res = await fetch(url_, { method: "PUT", credentials: "include" });
    await throwIfError(res);

    const obj: RoomDto = await res.json();
    return roomToModel(obj);
  };

  public startGame = async (roomId: UUID): Promise<Room> => {
    const url_ = this.url("rooms/:id/game", { params: { id: roomId } });
    const res = await fetch(url_, { method: "PUT", credentials: "include" });
    await throwIfError(res);

    const obj: RoomDto = await res.json();
    return roomToModel(obj);
  };
}
