import {Theme, defaultTheme} from "@/model/theme";
import Color from "color";
import { createContext, useEffect } from "react";
import { useCookies } from "react-cookie";

type SetTheme = (theme: Theme) => void;
type ThemeContextValue = {
  theme: Theme;
  setTheme: SetTheme;
};
export const ThemeContext = createContext<ThemeContextValue>(
  null as unknown as ThemeContextValue,
);

export const useTheme = (): ThemeContextValue => {
  const [cookies, setCookies, _] = useCookies(["theme"]);
  const theme: Theme = cookies["theme"] ?? defaultTheme;
  const setTheme = (theme: Theme) => {
    setCookies("theme", theme, { maxAge: 60 * 60 * 24 * 365 * 10 });
  };

  const applyTheme = () => {
    const root = document.documentElement;
    Object.entries(theme.colors).forEach(([key, value]) => {
      const colorKey = `--theme-${key}`;
      const colorValue = Color(value).array().join(" ");
      root.style.setProperty(colorKey, colorValue);
    });
  };
  useEffect(applyTheme, [theme]);
  return { theme, setTheme };
};
