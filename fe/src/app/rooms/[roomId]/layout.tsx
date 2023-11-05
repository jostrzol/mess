"use client";

import { useRoomWebsocket } from "@/api/room";
import {ConnectionStatus} from "@/components/connectionStatus";
import { RoomWsContext } from "@/contexts/roomWsContext";
import { UUID } from "crypto";

export type RoomPageParams = {
  params: {
    roomId: UUID;
  };
};

const RoomLayout = ({
  children,
  params,
}: RoomPageParams & { children: React.ReactNode }) => {
  const roomWsContextValue = useRoomWebsocket(params.roomId);
  return (
    <RoomWsContext.Provider value={roomWsContextValue}>
      <ConnectionStatus state={roomWsContextValue.readyState} />
      {children}
    </RoomWsContext.Provider>
  );
};

export default RoomLayout;
