import clsx from "clsx";
import { InputHTMLAttributes } from "react";
import {MdUploadFile} from "react-icons/md";

export const UploadFile = ({
  className,
  children,
  noIcon = false,
  ...props
}: InputHTMLAttributes<HTMLInputElement> & {noIcon?: boolean}) => {
  const { disabled } = props;
  return (
    <label
      className={clsx(
        "py-2",
        "px-3",
        "rounded-2xl",
        "bg-primary-dim/20",
        !disabled && "hover:bg-primary-dim/40",
        "active:bg-primary-dim/60",
        "focus:outline-none",
        "focus:ring",
        "focus:ring-primary",
        disabled && "opacity-40",
        "cursor-pointer",
        "group",
        "flex",
        "items-center",
        "gap-2",
      )}
    >
      <input type="file" className="hidden" accept=".hcl" {...props} />
      {!noIcon && <MdUploadFile />}
      {children}
    </label>
  );
};
