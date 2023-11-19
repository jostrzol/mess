import clsx from "clsx";
import { InputHTMLAttributes } from "react";

export const Input = ({
  className,
  children,
  ...props
}: InputHTMLAttributes<HTMLInputElement>) => {
  const { required, value, type = "text" } = props;
  const invalid = required && !value;
  const typeClassName: any = {
    text: invalid
      ? [
          "shadow-[0_2px_0_0_rgb(var(--theme-danger))]",
          "focus:shadow-[0_4px_0_0_rgb(var(--theme-danger))]",
        ]
      : [
          "shadow-[0_2px_0_0_rgb(var(--theme-primary))]",
          "focus:shadow-[0_4px_0_0_rgb(var(--theme-primary))]",
        ],
  }[type as string];
  return (
    <input
      className={clsx(
        "py-2",
        "px-3",
        "bg-transparent",
        "outline-none",
        "box-border",
        typeClassName,
        className,
      )}
      {...props}
    >
      {children}
    </input>
  );
};
