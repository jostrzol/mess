"use client";

import { useCookies } from "react-cookie";
import { defaultTheme, Theme, ThemeContext, useTheme } from "@/contexts/themeContext";
import {GamePage} from "@/pages/game";

const Home = () => {
  const [cookies, setCookies, _] = useCookies(["theme"]);
  const [theme, setTheme] = useTheme(cookies["theme"] ?? defaultTheme);
  const onThemeChange = (theme: Theme) => {
    setCookies("theme", theme, { maxAge: 60 * 60 * 24 * 365 * 10 });
    setTheme(theme);
  };
  return (
    <ThemeContext.Provider value={theme}>
      <body className="bg-background text-txt">
        <GamePage onThemeChange={onThemeChange}/>
      </body>
    </ThemeContext.Provider>
  );
};

export default Home;
