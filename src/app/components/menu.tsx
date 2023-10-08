import clsx from "clsx";
import { ReactNode, useContext, useState } from "react";
import { MdOutlineArrowForwardIos } from "react-icons/md";
import { Theme, ThemeContext, themes } from "../contexts/themeContext";

export interface MenuProps {
  onThemeChange?: (theme: Theme) => void;
}
export const Menu = ({ onThemeChange }: MenuProps) => {
  const [isMenuExpanded, setIsMenuExpanded] = useState(false);
  const currentTheme = useContext(ThemeContext);
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
        <div className="m-5">
          <MenuSection title="Theme">
            {Object.entries(themes).map(([name, colors]) => {
              const isSelected = name == currentTheme.name;
              return (
                <button
                  onClick={() => onThemeChange?.({ name, colors })}
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
                      "translate-y-[2px]",
                      "transition-transform",
                      isSelected && "scale-125",
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
