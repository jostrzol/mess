import { useRoomWebsocket } from "@/api/room";
import { createContext } from "react";

export type RoomWsContextValue = ReturnType<typeof useRoomWebsocket>;
export const RoomWsContext = createContext<RoomWsContextValue>(null!);
