"use client";

import { useCookies } from "react-cookie";
import { Board } from "./components/game/board";
import { Menu } from "./components/menu";
import { Theme, ThemeContext, themes, useTheme } from "./contexts/themeContext";

const Home = () => {
  const [cookies, setCookies, _] = useCookies(["theme"]);
  const [theme, setTheme] = useTheme(cookies["theme"] ?? themes["light"]);
  const onThemeChange = (theme: Theme) => {
    setCookies("theme", theme, { maxAge: 60 * 60 * 24 * 365 * 10 });
    setTheme(theme);
  };
  return (
    <ThemeContext.Provider value={theme}>
      <body className="bg-background text-txt">
        <Menu onThemeChange={onThemeChange} />
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
