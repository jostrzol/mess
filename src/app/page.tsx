"use client";

import clsx from "clsx";
import { Board } from "./components/game/board";
import { ThemeContext, themes, useTheme } from "./contexts/themeContext";
import { MdOutlineArrowForwardIos } from "react-icons/md";
import { useState } from "react";

const Home = () => {
  const [theme, setTheme] = useTheme(themes["light"]);
  const [isMenuExpanded, setIsMenuExpanded] = useState(false);
  return (
    <ThemeContext.Provider value={theme}>
      <body className="bg-background text-txt">
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
                  onClick={() => setTheme(theme)}
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
        <main className="h-screen flex justify-center items-center">
          <Board
            pieces={[]}
            board={{
              height: 8,
              width: 8,
            }}
          />
        </main>
      </body>
    </ThemeContext.Provider>
  );
};

export default Home;
