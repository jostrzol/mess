import { Event } from "@/api/schema/event";
import { url } from "@/api/url";
import { UUID } from "crypto";
import { ReactNode, createContext, useContext, useEffect } from "react";
import useWebSocket, { ReadyState, SendMessage } from "react-use-websocket";
import {
  SendJsonMessage,
  WebSocketLike,
} from "react-use-websocket/dist/lib/types";

export const RoomWebsocketContext = createContext<RoomWebsocketContextValue>(
  null!,
);
export type RoomWebsocketContextValue = {
  lastEvent: Event | null;
  sendMessage: SendMessage;
  sendJsonMessage: SendJsonMessage;
  readyState: ReadyState;
  getWebSocket: () => WebSocketLike | null;
};

export const useRoomWebsocket = <T extends Event>(handler?: {
  type: T["EventType"];
  onEvent: (event: T) => void;
}) => {
  const ws = useContext(RoomWebsocketContext);
  useEffect(() => {
    if (handler && ws.lastEvent?.EventType === handler.type) {
      handler.onEvent(ws.lastEvent as T);
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [ws.lastEvent]);
  return ws;
};

export const RoomWebsocketProvider = ({
  roomId,
  children,
}: {
  roomId: UUID;
  children?: ReactNode;
}) => {
  const url_ = url("rooms/:id/websocket", {
    params: { id: roomId },
    schema: "ws",
  });
  const { lastJsonMessage, ...rest } = useWebSocket(url_);
  const lastEvent = lastJsonMessage as Event | null;
  const value = { ...rest, lastEvent };
  return (
    <RoomWebsocketContext.Provider value={value}>
      {children}
    </RoomWebsocketContext.Provider>
  );
};
