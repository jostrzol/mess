"use client";

import clsx from "clsx";
import { ThemeContext, themes, useTheme } from "./contexts/themeContext";

export default function Home() {
  const [theme, setTheme] = useTheme(themes["light"]);
  return (
    <ThemeContext.Provider value={theme}>
      <body className="bg-background flex min-h-screen text-txt">
        <aside className="flex-1 flex flex-col justify-center items-stretch">
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
        </aside>
        <main
          className={clsx(
            "flex-[3]",
            "flex",
            "flex-col",
            "items-center",
            "justify-between",
            "p-24",
            "[&>*]:p-4",
            "[&>*]:m-4",
          )}
        >
          <div className={`bg-primary`}>Primary</div>
          <div className={`bg-primary-dim`}>Primary dim</div>
          <div className={`bg-txt`}>Text</div>
          <div className={`bg-txt-dim`}>Text dim</div>
          <div className={`bg-player-white`}>Player white</div>
          <div className={`bg-player-black`}>Player black</div>
        </main>
      </body>
    </ThemeContext.Provider>
  );
}
