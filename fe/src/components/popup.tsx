import clsx from "clsx";
import { ReactNode } from "react";

export const Popup = ({ children, className }: { children?: ReactNode, className?: string }) => {
  return (
    <>
      <div
        className={clsx(
          "fixed",
          "w-full",
          "h-full",
          "bg-background/40",
          "pointer-events-none",
        )}
      ></div>
      <div
        className={clsx(
          "fixed",
          "p-4",
          "bg-background",
          "rounded-2xl",
          "border-2",
          "border-primary",
          "flex",
          "flex-col",
          "items-stretch",
          className,
        )}
      >
        {children}
      </div>
    </>
  );
};
