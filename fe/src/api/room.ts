import { Room } from "@/model/room";
import { UUID } from "crypto";
import { MessApi } from "./messApi";
import { RoomDto, roomToModel } from "./schema/room";

export class RoomApi extends MessApi {
  public createRoom = async (): Promise<Room> => {
    const res = await this.fetch("rooms", {
      method: "POST",
      credentials: "include",
    });

    const obj: RoomDto = await res.json();
    return roomToModel(obj);
  };

  public joinRoom = async (id: UUID): Promise<Room> => {
    const res = await this.fetch("rooms/:id/players", {
      method: "PUT",
      params: { id: id },
      credentials: "include",
    });

    const obj: RoomDto = await res.json();
    return roomToModel(obj);
  };

  public getRoom = async (id: UUID): Promise<Room> => {
    const res = await this.fetch("rooms/:id", {
      method: "GET",
      params: { id },
      credentials: "include",
    });

    const obj: RoomDto = await res.json();
    return roomToModel(obj);
  };

  public getRules = async (roomId: UUID): Promise<string> => {
    const res = await this.fetch("rooms/:id/rules", {
      method: "GET",
      params: { id: roomId },
      credentials: "include",
    });

    return await res.text();
  };

  public startGame = async (roomId: UUID): Promise<Room> => {
    const res = await this.fetch("rooms/:id/game", {
      method: "PUT",
      params: { id: roomId },
      credentials: "include",
    });

    const obj: RoomDto = await res.json();
    return roomToModel(obj);
  };

  public saveRules = async (
    roomId: UUID,
    filename: string,
    data: string,
  ): Promise<void> => {
    await this.fetch("rooms/:id/rules/:filename", {
      method: "PUT",
      params: { id: roomId, filename },
      credentials: "include",
      body: data,
    });
  };
}
