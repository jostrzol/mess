"use client";

import { themes } from "@/model/theme";
import clsx from "clsx";
import { ReactNode, useContext, useState } from "react";
import { MdOutlineArrowForwardIos } from "react-icons/md";
import { ThemeContext } from "../contexts/themeContext";
import { Logo } from "./logo";

export const Menu = () => {
  const [isMenuExpanded, setIsMenuExpanded] = useState(false);
  const { theme, setTheme } = useContext(ThemeContext);
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
        "z-10",
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
          "overflow-hidden",
          "bg-background/80",
          "whitespace-nowrap",
        )}
      >
        <div className="m-5 flex flex-col gap-4">
          <section
            className={clsx(
              "flex",
              "m-2",
              "py-2",
              "px-3",
              "bg-primary-dim/20",
              "rounded-2xl",
            )}
          >
            <Logo size={50} className={clsx("min-w-[50px]", "mr-2")} />
            <h1>mess</h1>
          </section>
          <MenuSection title="Theme">
            {Object.entries(themes).map(([name, colors]) => {
              const isSelected = name == theme.name;
              return (
                <button
                  onClick={() => setTheme({ name, colors })}
                  className={clsx(
                    "mb-2",
                    "bg-background",
                    "border-primary",
                    "max-w-fit",
                  )}
                  key={name}
                >
                  <span
                    className={clsx(
                      "w-3",
                      "h-3",
                      "inline-block",
                      "mr-2",
                      "rounded-full",
                      "border-2",
                      "translate-y-[2px]",
                      "transition-transform",
                      isSelected && "scale-150",
                    )}
                    style={{
                      backgroundColor: colors.background,
                      borderColor: colors.primary,
                    }}
                  />
                  <span className={clsx(isSelected || "text-txt-dim")}>
                    {name[0].toUpperCase() + name.slice(1).toLowerCase()}
                  </span>
                </button>
              );
            })}
          </MenuSection>
        </div>
      </div>
      <button
        className={clsx(
          "flex-initial",
          "shrink-0",
          "p-2",
          "h-full",
          "bg-background/80",
          "group",
          "border-r-2",
          "border-txt-dim",
          "border-dashed",
          "outline-none",
        )}
        onClick={() => setIsMenuExpanded(!isMenuExpanded)}
      >
        <div
          className={clsx(
            "transition-transform",
            "duration-300",
            "delay-300",
            isMenuExpanded && "rotate-180",
          )}
        >
          <div
            className={clsx(
              "absolute",
              "group-focus:animate-ping-1",
              "group-hover:animate-ping-1",
            )}
          >
            <MdOutlineArrowForwardIos />
          </div>
          <MdOutlineArrowForwardIos />
        </div>
      </button>
    </aside>
  );
};

interface MenuSectionProps {
  title: string;
  children: ReactNode;
}
const MenuSection = ({ title, children }: MenuSectionProps) => {
  return (
    <section>
      <h1 className={clsx("mb-4", "font-semibold")}>{title}</h1>
      <div className={clsx("flex", "flex-col", "ml-2")}>{children}</div>
    </section>
  );
};
