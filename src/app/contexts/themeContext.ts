import {
  Dispatch,
  SetStateAction,
  createContext,
  useEffect,
  useState,
} from "react";
import Color from "color";

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
    background: "#ffffff",
    primary: "#1e40af",
    "primary-dim": "#2563eb",
    txt: "#292524",
    "txt-dim": "#57534e",
    "player-white": "#ffffff",
    "player-black": "#000000",
  },
  dark: {
    background: "#000000",
    primary: "#1e40af",
    "primary-dim": "#2563eb",
    txt: "#78716c",
    "txt-dim": "#44403c",
    "player-white": "#ffffff",
    "player-black": "#000000",
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
