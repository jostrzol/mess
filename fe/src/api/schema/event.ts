export interface RoomChanged {
  EventType: "RoomChanged";
}

export interface GameStarted {
  EventType: "GameStarted";
}

export type Event = RoomChanged | GameStarted;
