"use client";

import { RoomWebsocketProvider } from "@/contexts/roomWsContext";
import { RuleFilesProvider } from "@/contexts/rulesContext";
import { UUID } from "crypto";
import { ReactNode } from "react";

export type RoomPageParams = {
  params: {
    roomId: UUID;
  };
};

const RoomLayout = ({ children }: { children: ReactNode }) => {
  return (
    <RoomWebsocketProvider>
      <RuleFilesProvider>{children}</RuleFilesProvider>
    </RoomWebsocketProvider>
  );
};

export default RoomLayout;
