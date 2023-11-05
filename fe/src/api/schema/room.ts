import {Room} from "@/model/room";
import {UUID} from "crypto";

export interface RoomDto {
  ID: UUID;
  Players: number;
  PlayersNeeded: number;
  IsStartable: boolean;
  IsStarted: boolean;
}

export const roomToModel = (room: RoomDto): Room => {
  return {
    id: room.ID,
    players: room.Players,
    playersNeeded: room.PlayersNeeded,
    isStartable: room.IsStartable,
    isStarted: room.IsStarted,
  };
};

