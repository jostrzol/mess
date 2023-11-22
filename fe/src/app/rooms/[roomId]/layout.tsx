"use client";

import { RoomWebsocketProvider } from "@/contexts/roomWsContext";
import { UUID } from "crypto";
import { ReactNode } from "react";

export type RoomPageParams = {
  params: {
    roomId: UUID;
  };
};

const RoomLayout = ({ children }: { children: ReactNode }) => {
  return <RoomWebsocketProvider>{children}</RoomWebsocketProvider>;
};

export default RoomLayout;
