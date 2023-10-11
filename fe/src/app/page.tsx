"use client";

import { useCookies } from "react-cookie";
import { Board } from "./components/game/board";
import { Menu } from "./components/menu";
import { Theme, ThemeContext, themes, useTheme } from "./contexts/themeContext";
import { Piece as PieceConfig } from "./model/piece";
import { PieceType } from "./model/pieceType";
import { Square } from "./model/square";

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
    queen: {
      code: "queen",
      name: "Queen",
      iconUri: "/pieces/queen.svg",
    },
    rook: {
      code: "rook",
      name: "Rook",
      iconUri: "/pieces/rook.svg",
    },
    bishop: {
      code: "bishop",
      name: "Bishop",
      iconUri: "/pieces/bishop.svg",
    },
    knight: {
      code: "knight",
      name: "Knight",
      iconUri: "/pieces/knight.svg",
    },
    pawn: {
      code: "pawn",
      name: "Pawn",
      iconUri: "/pieces/pawn.svg",
    },
  } satisfies Record<string, PieceType>;
  const makeMoves = (square: Square) => {
    const offsets = [
      [1, 1],
      [0, 1],
      [-1, 1],
      [1, -1],
      [0, -1],
      [-1, -1],
      [-1, 0],
      [1, 0],
    ];
    return offsets
      .map((offset): Square => [square[0] + offset[0], square[1] + offset[1]])
      .filter((dest) => dest[0] >= 0 && dest[1] >= 0);
  };
  const pieces: PieceConfig[] = Object.values(pieceTypes).flatMap((type, i) => [
    {
      square: [3, i],
      color: "black",
      type,
      validMoves: makeMoves([3, i]),
    },
    {
      square: [4, i],
      color: "black",
      type,
      validMoves: makeMoves([4, i]),
    },
    {
      square: [5, i],
      color: "white",
      type,
      validMoves: makeMoves([5, i]),
    },
    {
      square: [6, i],
      color: "white",
      type,
      validMoves: makeMoves([6, i]),
    },
  ]);
  return (
    <ThemeContext.Provider value={theme}>
      <body className="bg-background text-txt">
        <Menu onThemeChange={onThemeChange} />
        <main className="h-screen flex justify-center items-center">
          <Board
            pieces={pieces}
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
