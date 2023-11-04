import { UUID } from "crypto";

export interface Room {
  id: UUID;
  players: number;
  playersNeeded: number;
  isStarted: boolean;
  isStartable: boolean;
}
