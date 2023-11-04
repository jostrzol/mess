import { ThemeContext } from "@/contexts/themeContext";
import clsx from "clsx";
import { useContext } from "react";
import { ReadyState } from "react-use-websocket";

export const ConnectionStatus = ({
  state,
  className,
}: {
  state: ReadyState;
  className?: string;
}) => {
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
  }[state];
  return (
    <div
      className={clsx(
        "w-4",
        "h-4",
        "rounded-full",
        "has-tooltip",
        "absolute",
        "top-5",
        "right-5",
        className,
      )}
      style={{ backgroundColor: color }}
    >
      <span
        className={clsx(
          "tooltip",
          "rounded",
          "shadow-lg",
          "p-1",
          "bg-primary-dim/50",
          "-left-2",
          "-translate-y-2",
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
