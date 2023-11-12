import clsx from "clsx";
import { ReactNode } from "react";

export const Window = ({
  opaque = false,
  title,
  className,
  children,
}: {
  opaque?: boolean;
  title?: string;
  className?: string;
  children?: ReactNode;
}) => {
  return (
    <>
      <div
        className={clsx(
          "p-4",
          opaque ? "bg-background/90" : "bg-background",
          "rounded-2xl",
          "border-2",
          "border-primary",
          "flex",
          "flex-col",
          "items-stretch",
          className,
        )}
      >
        {title && <h2 className="font-semibold text-center">{title}</h2>}
        {children}
      </div>
    </>
  );
};
