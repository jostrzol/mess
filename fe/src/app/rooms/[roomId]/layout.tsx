"use client";

import { ConnectionStatus } from "@/components/connectionStatus";
import {
  RoomWebsocketProvider,
  useRoomWebsocket,
} from "@/contexts/roomWsContext";
import {useIsFetching} from "@tanstack/react-query";
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
      <RoomLayoutInner>{children}</RoomLayoutInner>
    </RoomWebsocketProvider>
  );
};

const RoomLayoutInner = ({ children }: { children: ReactNode }) => {
  const { readyState } = useRoomWebsocket();
  const fetching = useIsFetching()
  return (
    <>
      <ConnectionStatus websocketStatus={readyState} isFetching={fetching > 0} />
      {children}
    </>
  );
};

export default RoomLayout;
