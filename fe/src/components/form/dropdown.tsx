import clsx from "clsx";
import { ReactNode, useState } from "react";
import ReactDropdown, { ReactDropdownProps } from "react-dropdown";
import { MdArrowDropDown, MdArrowDropUp } from "react-icons/md";

export type DropdownPros = ReactDropdownProps & {
  arrowFloating?: boolean;
};

export const Dropdown = ({
  arrowFloating = false,
  className,
  menuClassName,
  controlClassName,
  arrowClassName,
  arrowClosed,
  arrowOpen,
  ...props
}: DropdownPros) => {
  const [isFocused, setIsFocused] = useState(false);
  return (
    <ReactDropdown
      onFocus={(isNotFocused) => setIsFocused(!isNotFocused)}
      onChange={() => setIsFocused(false)}
      className={clsx("relative cursor-pointer pl-2", isFocused && "border-l-2", className)}
      controlClassName={clsx("relative flex items-center gap-2", controlClassName)}
      menuClassName={clsx("absolute bg-background/90 min-w-full hover:[&>*]:translate-x-1 pl-2 -left-[2px] border-l-2", menuClassName)}
      arrowClosed={
        <Arrow floating={arrowFloating} className={arrowClassName}>
          {arrowClosed || <MdArrowDropDown />}
        </Arrow>
      }
      arrowOpen={
        <Arrow floating={arrowFloating} className={arrowClassName}>
          {arrowOpen || <MdArrowDropUp />}
        </Arrow>
      }
      {...props}
    />
  );
};

const Arrow = ({
  floating = false,
  className,
  children,
}: {
  floating?: boolean;
  className?: string;
  children: ReactNode;
}) => (
  <span
    className={clsx(
      floating && "absolute top-1/2 left-full -translate-y-1/2",
      className,
    )}
  >
    {children}
  </span>
);
