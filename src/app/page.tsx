"use client";

import {Board} from "./components/game/board";
import {Menu} from "./components/menu";
import {ThemeContext, themes, useTheme} from "./contexts/themeContext";

const Home = () => {
  const [theme, setTheme] = useTheme(themes["light"]);
  return (
    <ThemeContext.Provider value={theme}>
      <body className="bg-background text-txt">
        <Menu onThemeChange={(theme) => setTheme(theme)} />
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
