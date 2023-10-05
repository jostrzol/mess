import {Dispatch, SetStateAction, createContext, useEffect, useState} from "react";
import * as colors from "tailwindcss/colors";

export interface Theme {
  "background": string,
  "primary": string,
  "primary-dim": string,
  "txt": string,
  "txt-dim": string,
  "player-white": string,
  "player-black": string,
};

export const themes = {
  "light": {
    "background": colors.white,
    "primary": colors.blue[600],
    "primary-dim": colors.blue[800],
    "txt": colors.stone[800],
    "txt-dim": colors.stone[600],
    "player-white": colors.white,
    "player-black": colors.black,
  },
  "dark": {
    "background": colors.black,
    "primary": colors.blue[600],
    "primary-dim": colors.blue[800],
    "txt": colors.stone[500],
    "txt-dim": colors.stone[700],
    "player-white": colors.white,
    "player-black": colors.black,
  },
} satisfies Record<string, Theme>;

export const ThemeContext = createContext<Theme>(themes["light"]);

export const useTheme = (initialTheme: Theme): [Theme, Dispatch<SetStateAction<Theme>>] => {
  const [theme, setTheme] = useState<Theme>(initialTheme);
  const applyTheme = () => {
    const root = document.documentElement;
    Object.entries(theme).forEach(([key, value]) => {
      const colorKey = `--theme-${key}`;
      console.log(`${colorKey}: ${value}`);
      root.style.setProperty(colorKey, value);
    });
  };
  useEffect(applyTheme, [theme]);
  return [theme, setTheme];
};
