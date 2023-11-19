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
    const url_ = this.url("rooms/:id/players", { params: { id } });
    const res = await fetch(url_, { method: "PUT", credentials: "include" });
    await throwIfError(res);

    const obj: RoomDto = await res.json();
    return roomToModel(obj);
  };

  public getRoom = async (id: UUID): Promise<Room> => {
    const url_ = this.url("rooms/:id", { params: { id } });
    const res = await fetch(url_, { method: "GET", credentials: "include" });
    await throwIfError(res);

    const obj: RoomDto = await res.json();
    return roomToModel(obj);
  };

  public getRules = async (roomId: UUID): Promise<string> => {
    const url_ = this.url("rooms/:id/rules", { params: { id: roomId } });
    const res = await fetch(url_, { method: "GET", credentials: "include" });
    await throwIfError(res);

    return await res.text();
  };

  public startGame = async (roomId: UUID): Promise<Room> => {
    const url_ = this.url("rooms/:id/game", { params: { id: roomId } });
    const res = await fetch(url_, { method: "PUT", credentials: "include" });
    await throwIfError(res);

    const obj: RoomDto = await res.json();
    return roomToModel(obj);
  };

  public saveRules = async (roomId: UUID, filename: string, data: string): Promise<void> => {
    const url_ = this.url("rooms/:id/rules/:filename", { params: { id: roomId, filename } });
    const res = await fetch(url_, { method: "PUT", credentials: "include", body: data });
    await throwIfError(res);
  };
}
