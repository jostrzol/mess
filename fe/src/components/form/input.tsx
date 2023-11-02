import clsx from "clsx";
import { InputHTMLAttributes } from "react";

export const Input = ({
  className,
  children,
  ...props
}: InputHTMLAttributes<HTMLInputElement>) => {
  return (
    <input
      className={clsx(
        "m-2",
        "py-2",
        "px-3",
        "bg-transparent",
        "shadow-[0_2px_0_0_rgb(var(--theme-primary))]",
        "focus:shadow-[0_4px_0_0_rgb(var(--theme-primary))]",
        "outline-none",
        "box-border",
        className,
      )}
      {...props}
    >
      {children}
    </input>
  );
};
