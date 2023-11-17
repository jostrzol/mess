import colors from "tailwindcss/colors";

export interface ThemeColors {
  background: string;
  primary: string;
  "primary-dim": string;
  txt: string;
  "txt-dim": string;
  "player-white": string;
  "player-black": string;
  danger: string;
  "danger-dim": string;
  warn: string;
  "warn-dim": string;
  "success-strong": string;
  success: string;
  "success-dim": string;
}

export interface Theme {
  name: string;
  colors: ThemeColors;
}

export const themes = {
  light: {
    background: colors.slate[50],
    primary: colors.slate[700],
    "primary-dim": colors.slate[500],
    txt: colors.slate[900],
    "txt-dim": colors.slate[400],
    "player-white": colors.slate[50],
    "player-black": colors.slate[950],
    danger: colors.rose[600],
    "danger-dim": colors.rose[400],
    warn: colors.amber[500],
    "warn-dim": colors.amber[300],
    "success-strong": colors.lime[600],
    success: colors.lime[400],
    "success-dim": colors.lime[200],
  },
  dark: {
    background: colors.slate[950],
    primary: colors.slate[700],
    "primary-dim": colors.slate[500],
    txt: colors.slate[400],
    "txt-dim": colors.slate[600],
    "player-white": colors.slate[50],
    "player-black": colors.slate[950],
    danger: colors.rose[600],
    "danger-dim": colors.rose[800],
    warn: colors.amber[500],
    "warn-dim": colors.amber[300],
    "success-strong": colors.lime[600],
    success: colors.lime[400],
    "success-dim": colors.lime[200],
  },
} satisfies Record<string, ThemeColors>;

export const defaultTheme: Theme = {
  name: "light",
  colors: themes["light"],
};
