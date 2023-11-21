import clsx from "clsx";
import { ReactNode } from "react";
import { Window } from "./window";

export const Popup = ({
  children,
  className,
}: {
  children?: ReactNode;
  className?: string;
}) => {
  return (
    <>
      <div
        className={clsx(
          "fixed",
          "z-40",
          "w-full",
          "h-full",
          "bg-background/40",
          "pointer-events-none",
        )}
      ></div>
      <Window className={clsx("fixed", "z-50", className)}>{children}</Window>
    </>
  );
};
