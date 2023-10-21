import Color from "color";
import {
  Dispatch,
  SetStateAction,
  createContext,
  useEffect,
  useState,
} from "react";
import * as colors from "tailwindcss/colors";

export interface ThemeColors {
  background: string;
  primary: string;
  "primary-dim": string;
  txt: string;
  "txt-dim": string;
  "player-white": string;
  "player-black": string;
  "danger": string;
  "danger-dim": string;
}

export interface Theme {
  name: string;
  colors: ThemeColors;
}

export const themes = {
  light: {
    background: colors.slate[50],
    primary: colors.blue[700],
    "primary-dim": colors.blue[500],
    txt: colors.slate[900],
    "txt-dim": colors.slate[400],
    "player-white": colors.slate[50],
    "player-black": colors.slate[950],
    danger: colors.rose[600],
    "danger-dim": colors.rose[400],
  },
  dark: {
    background: colors.slate[950],
    primary: colors.blue[700],
    "primary-dim": colors.blue[500],
    txt: colors.slate[400],
    "txt-dim": colors.slate[600],
    "player-white": colors.slate[50],
    "player-black": colors.slate[950],
    danger: colors.rose[600],
    "danger-dim": colors.rose[800],
  },
} satisfies Record<string, ThemeColors>;

export const defaultTheme: Theme = {
  name: "light",
  colors: themes["light"],
};
export const ThemeContext = createContext<Theme>(defaultTheme);

type SetTheme = Dispatch<SetStateAction<Theme>>;

export const useTheme = (initialTheme: Theme): [Theme, SetTheme] => {
  const [theme, setTheme] = useState<Theme>(initialTheme);
  const applyTheme = () => {
    console.log(theme.colors)
    const root = document.documentElement;
    Object.entries(theme.colors).forEach(([key, value]) => {
      const colorKey = `--theme-${key}`;
      const colorValue = Color(value).array().join(" ");
      root.style.setProperty(colorKey, colorValue);
    });
  };
  useEffect(applyTheme, [theme]);
  return [theme, setTheme];
};