"use client";

import {ThemeContext, themes, useTheme} from "./contexts/themeContext";

export default function Home() {
  const [theme, setTheme] = useTheme(themes["light"]);
  return (
    <ThemeContext.Provider value={theme}>
      <main className="bg-background flex min-h-screen flex-col items-center justify-between p-24">
        <div className={`bg-primary p-2`}>
          Primary
        </div>
        <div className={`bg-primary-dim p-2`}>
          Primary dim
        </div>
        <div className={`bg-txt p-2`}>
          Text
        </div>
        <div className={`bg-txt-dim p-2`}>
          Text dim
        </div>
        <div className={`bg-player-white p-2`}>
          Player white
        </div>
        <div className={`bg-player-black p-2`}>
          Player black
        </div>
      </main>
    </ThemeContext.Provider>
  );
}
