"use client";

import { useCookies } from "react-cookie";
import { Board } from "./components/game/board";
import { Menu } from "./components/menu";
import { Theme, ThemeContext, themes, useTheme } from "./contexts/themeContext";
import { PieceType } from "./model/pieceType";
import { Piece } from "./components/game/piece";

const Home = () => {
  const [cookies, setCookies, _] = useCookies(["theme"]);
  const [theme, setTheme] = useTheme(cookies["theme"] ?? themes["light"]);
  const onThemeChange = (theme: Theme) => {
    setCookies("theme", theme, { maxAge: 60 * 60 * 24 * 365 * 10 });
    setTheme(theme);
  };
  const pieceTypes = {
    king: {
      code: "king",
      name: "King",
      iconUri: "/pieces/king.svg",
    },
  } satisfies Record<string, PieceType>;
  return (
    <ThemeContext.Provider value={theme}>
      <body className="bg-background text-txt">
        <Menu onThemeChange={onThemeChange} />
        <main className="h-screen flex justify-center items-center">
          <Board
            pieces={[
              {
                location: [6, 5],
                color: "black",
                type: pieceTypes.king,
              },
              {
                location: [6, 6],
                color: "black",
                type: pieceTypes.king,
              },
            ]}
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
