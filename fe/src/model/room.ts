import {UUID} from "crypto";

export const MaxPlayers = 2;

export class Room {
  id: UUID;
  players: number;

  constructor(id: UUID, players: number) {
    this.id = id
    this.players = players
  }

  isReady(): boolean {
    return this.players === MaxPlayers
  }
}

