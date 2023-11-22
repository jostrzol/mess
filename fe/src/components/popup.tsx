import clsx from "clsx";
import { ReactNode } from "react";
import { Window } from "./window";

type Position = "middle" | "bottom";

export const Popup = ({
  title,
  position = "middle",
  modal = false,
  children,
  className,
}: {
  title?: string;
  position?: Position;
  modal?: boolean;
  children?: ReactNode;
  className?: string;
}) => {
  return (
    <>
      {modal && (
        <div
          className={clsx(
            "fixed",
            "z-40",
            "w-full",
            "h-full",
            "bg-background/40",
          )}
        />
      )}
      <Window
        title={title}
        opaque={!modal}
        className={clsx(
          "fixed",
          "z-50",
          position === "bottom" && "bottom-0 m-4 max-w-[90%]",
          position === "middle" &&
            "top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2",
          className,
        )}
      >
        {children}
      </Window>
    </>
  );
};
