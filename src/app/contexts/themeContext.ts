import {
  Dispatch,
  SetStateAction,
  createContext,
  useEffect,
  useState,
} from "react";

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

export const useTheme = (
  initialTheme: Theme,
): [Theme, Dispatch<SetStateAction<Theme>>] => {
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
