import { useRoomWebsocket } from "@/contexts/roomWsContext";
import { ThemeContext } from "@/contexts/themeContext";
import { useIsFetching } from "@tanstack/react-query";
import clsx from "clsx";
import { useContext } from "react";
import { ReadyState } from "react-use-websocket";

export const ConnectionStatus = ({ className }: { className?: string }) => {
  const { readyState: websocketStatus } = useRoomWebsocket();
  const fetching = useIsFetching();
  const {
    theme: { colors },
  } = useContext(ThemeContext);
  const [tooltip, color] = {
    [ReadyState.CONNECTING]: ["Connecting", colors.warn],
    [ReadyState.OPEN]: ["Connection open", colors.success],
    [ReadyState.CLOSING]: ["Closing connection", colors.warn],
    [ReadyState.CLOSED]: ["Connection closed", colors.danger],
    [ReadyState.UNINSTANTIATED]: [
      "Connection uninstantiated",
      colors["txt-dim"],
    ],
  }[websocketStatus];
  return (
    <div className={clsx("relative", "has-tooltip", className)}>
      <div
        className={clsx("w-4", "h-4", "rounded-full", "has-tooltip")}
        style={{ backgroundColor: color }}
      />
      <div
        className={clsx(
          "absolute",
          "top-0",
          "left-0",
          "w-4",
          "h-4",
          "rounded-full",
          "has-tooltip",
          fetching > 0 && "animate-ping",
        )}
        style={{ backgroundColor: color }}
      />
      <span
        className={clsx(
          "tooltip",
          "rounded",
          "shadow-lg",
          "p-1",
          "bg-primary-dim/50",
          "-left-2",
          "-translate-y-3/4",
          "-translate-x-full",
          "whitespace-nowrap",
          "transition-opacity",
          "text-sm",
        )}
      >
        {tooltip}
      </span>
    </div>
  );
};
