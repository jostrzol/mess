"use client";

import { themes } from "@/model/theme";
import clsx from "clsx";
import { ReactNode, useContext } from "react";
import { MdOutlineArrowBackIos } from "react-icons/md";
import { ThemeContext } from "../contexts/themeContext";
import { Logo } from "./logo";

export const Menu = ({
  open = true,
  onClose,
  children,
}: {
  open?: boolean;
  onClose?: () => void;
  children?: ReactNode;
}) => {
  return (
    <>
      <aside
        className={clsx(
          "fixed top-0 left-0",
          "max-w-full w-96 h-screen",
          "flex",
          "z-50",
          "transition-transform duration-500",
          !open && "-translate-x-full",
          "overflow-hidden",
          "bg-background/80",
          "border-r-2 border-dashed",
        )}
      >
        <div className="grow my-4 mr-2 ml-4 flex flex-col gap-4">
          <MenuLogo />
          {children}
          <ThemeMenuSection />
        </div>
        <button
          className={clsx(
            "flex-initial",
            "shrink-0",
            "p-2",
            "h-full",
            "bg-background/80",
            "group",
            "outline-none",
          )}
          onClick={() => onClose?.()}
        >
          <div
            className={clsx(
              "absolute",
              "group-focus:animate-ping-1",
              "group-hover:animate-ping-1",
            )}
          >
            <MdOutlineArrowBackIos />
          </div>
          <MdOutlineArrowBackIos />
        </button>
      </aside>
      <div
        className={clsx(
          "fixed z-40 top-0 left-0 h-screen w-screen",
          "bg-background/50",
          "transition-opacity duration-500",
          !open && ["opacity-0", "pointer-events-none"],
        )}
        onClick={() => onClose?.()}
      />
    </>
  );
};

const MenuLogo = () => (
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
);

const ThemeMenuSection = () => {
  const { theme: selectedTheme, setTheme } = useContext(ThemeContext);
  return (
    <MenuSection title="Theme">
      {Object.values(themes).map((theme) => {
        const isSelected = theme.name == selectedTheme.name;
        return (
          <button
            onClick={() => setTheme(theme)}
            className={clsx(
              "mb-2",
              "bg-background",
              "border-primary",
              "max-w-fit",
            )}
            key={theme.name}
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
                backgroundColor: theme.colors.background,
                borderColor: theme.colors.primary,
              }}
            />
            <span className={clsx(isSelected || "text-txt-dim")}>
              {theme.name[0].toUpperCase() + theme.name.slice(1).toLowerCase()}
            </span>
          </button>
        );
      })}
    </MenuSection>
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
