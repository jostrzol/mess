"use client";

import { ConnectionStatus } from "@/components/connectionStatus";
import {
  RoomWebsocketProvider,
  useRoomWebsocket,
} from "@/contexts/roomWsContext";
import { UUID } from "crypto";
import { ReactNode } from "react";

export type RoomPageParams = {
  params: {
    roomId: UUID;
  };
};

const RoomLayout = ({
  children,
  params,
}: RoomPageParams & { children: ReactNode }) => {
  return (
    <RoomWebsocketProvider roomId={params.roomId}>
      <RoomLayoutInner>{children}</RoomLayoutInner>
    </RoomWebsocketProvider>
  );
};

const RoomLayoutInner = ({ children }: { children: ReactNode }) => {
  const { readyState } = useRoomWebsocket();
  return (
    <>
      <ConnectionStatus state={readyState} />
      {children}
    </>
  );
};

export default RoomLayout;
