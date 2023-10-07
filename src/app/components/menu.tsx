import clsx from "clsx";
import { SetTheme, Theme, themes } from "../contexts/themeContext";
import { useState } from "react";
import {MdOutlineArrowForwardIos} from "react-icons/md";

export interface MenuProps {
  onThemeChange?: (theme: Theme) => void
}
export const Menu = ({onThemeChange}: MenuProps) => {
  const [isMenuExpanded, setIsMenuExpanded] = useState(false);
  return (
    <aside
      className={clsx(
        "flex",
        "justify-start",
        "sticky",
        "max-w-fit",
        "top-0",
        "left-0",
        "h-screen",
        "-mb-[100vh]",
      )}
    >
      <div
        className={clsx(
          "transition-[width]",
          "duration-500",
          isMenuExpanded && "w-96",
          "w-0",
          "flex",
          "flex-col",
          "justify-center",
          "items-stretch",
          "overflow-hidden",
          "bg-background/80",
        )}
      >
        {Object.entries(themes).map(([name, theme]) => {
          return (
            <button
              onClick={() => onThemeChange?.(theme)}
              className="p-2 m-4 border-primary border-2 rounded-md"
              key={name}
            >
              {name}
            </button>
          );
        })}
      </div>
      <button
        className={clsx(
          "flex-initial",
          "shrink-0",
          "p-2",
          "h-full",
          "bg-background/80",
          "group",
        )}
        onClick={() => setIsMenuExpanded(!isMenuExpanded)}
      >
        <div
          className={clsx(
            "transition",
            "duration-300",
            "delay-300",
            isMenuExpanded && "rotate-180",
          )}
        >
          <div className={clsx("absolute", "group-hover:animate-ping-1")}>
            <MdOutlineArrowForwardIos />
          </div>
          <MdOutlineArrowForwardIos />
        </div>
      </button>
    </aside>
  );
};
