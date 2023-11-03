import clsx from "clsx";
import { ButtonHTMLAttributes } from "react";

export const Button = ({
  className,
  children,
  ...props
}: ButtonHTMLAttributes<HTMLButtonElement>) => {
  return (
    <button
      className={clsx(
        "py-2",
        "px-3",
        "bg-primary-dim/20",
        "rounded-2xl",
        "enabled:hover:bg-primary-dim/40",
        "active:bg-primary-dim/60",
        "focus:outline-none",
        "focus:ring",
        "focus:ring-primary",
        "disabled:opacity-40",
        className,
      )}
      {...props}
    >
      {children}
    </button>
  );
};
