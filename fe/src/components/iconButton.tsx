import clsx from "clsx";
import { ButtonHTMLAttributes } from "react";

export const IconButton = ({
  children,
  className,
  ...props
}: ButtonHTMLAttributes<HTMLButtonElement>) => {
  return (
    <button
      className={clsx(
        "p-2",
        "self-center",
        "bg-primary-dim/20",
        "hover:bg-primary-dim/40",
        "cursor-pointer",
        "rounded-full",
        className,
      )}
      {...props}
    >
      {children}
    </button>
  );
};
