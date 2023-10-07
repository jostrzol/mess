import {
  Dispatch,
  SetStateAction,
  createContext,
  useEffect,
  useState,
} from "react";
import Color from "color";
import * as colors from "tailwindcss/colors";

export interface Theme {
  background: string;
  primary: string;
  "primary-dim": string;
  txt: string;
  "txt-dim": string;
  "player-white": string;
  "player-black": string;
}

export const themes = {
  light: {
    background: colors.slate[50],
    primary: "#1e40af",
    "primary-dim": "#2563eb",
    txt: colors.slate[900],
    "txt-dim": colors.slate[800],
    "player-white": colors.slate[50],
    "player-black": colors.slate[950],
  },
  dark: {
    background: colors.slate[950],
    primary: "#1e40af",
    "primary-dim": "#2563eb",
    txt: colors.slate[400],
    "txt-dim": colors.slate[500],
    "player-white": colors.slate[50],
    "player-black": colors.slate[950],
  },
} satisfies Record<string, Theme>;

export const ThemeContext = createContext<Theme>(themes["light"]);

export type SetTheme = Dispatch<SetStateAction<Theme>>;

export const useTheme = (initialTheme: Theme): [Theme, SetTheme] => {
  const [theme, setTheme] = useState<Theme>(initialTheme);
  const applyTheme = () => {
    const root = document.documentElement;
    Object.entries(theme).forEach(([key, value]) => {
      const colorKey = `--theme-${key}`;
      const colorValue = Color(value).array().join(" ");
      console.log(`${colorKey}: ${colorValue}`);
      root.style.setProperty(colorKey, colorValue);
    });
  };
  useEffect(applyTheme, [theme]);
  return [theme, setTheme];
};
