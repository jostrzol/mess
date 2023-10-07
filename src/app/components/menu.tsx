import clsx from "clsx";
import { Theme, themes } from "../contexts/themeContext";
import { ReactNode, useState } from "react";
import { MdOutlineArrowForwardIos } from "react-icons/md";

export interface MenuProps {
  onThemeChange?: (theme: Theme) => void;
}
export const Menu = ({ onThemeChange }: MenuProps) => {
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
          "overflow-hidden",
          "bg-background/80",
          "whitespace-nowrap",
        )}
      >
        <div className="m-5">
          <MenuSection title="Theme">
            {Object.entries(themes).map(([name, theme]) => {
              return (
                <button
                  onClick={() => onThemeChange?.(theme)}
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
                      "w-4",
                      "h-4",
                      "inline-block",
                      "mr-2",
                      "rounded-full",
                      "border-2",
                      "translate-y-[1px]",
                    )}
                    style={{
                      backgroundColor: theme.background,
                      borderColor: theme.primary,
                    }}
                  />
                  {name[0].toUpperCase() + name.slice(1).toLowerCase()}
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
