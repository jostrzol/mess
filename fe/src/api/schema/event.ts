export interface RoomChanged {
  EventType: "RoomChanged";
}

export interface GameChanged {
  EventType: "GameChanged";
}

export type Event = RoomChanged | GameChanged;
