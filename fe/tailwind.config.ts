import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    colors: {
      transparent: "transparent",
      background: "rgb(var(--theme-background) / <alpha-value>)",
      primary: {
        DEFAULT: "rgb(var(--theme-primary) / <alpha-value>)",
        dim: "rgb(var(--theme-primary-dim) / <alpha-value>)",
      },
      txt: {
        DEFAULT: "rgb(var(--theme-txt) / <alpha-value>)",
        dim: "rgb(var(--theme-txt-dim) / <alpha-value>)",
      },
      player: {
        white: "rgb(var(--theme-player-white) / <alpha-value>)",
        DEFAULT: "rgb(var(--player-color) / <alpha-value>)",
        black: "rgb(var(--theme-player-black) / <alpha-value>)",
      },
      opponent: "rgb(var(--opponent-color) / <alpha-value>)",
      danger: {
        DEFAULT: "rgb(var(--theme-danger) / <alpha-value>)",
        dim: "rgb(var(--theme-danger-dim) / <alpha-value>)",
      },
      warn: {
        DEFAULT: "rgb(var(--theme-warn) / <alpha-value>)",
        dim: "rgb(var(--theme-warn-dim) / <alpha-value>)",
      },
      success: {
        strong: "rgb(var(--theme-success-strong) / <alpha-value>)",
        DEFAULT: "rgb(var(--theme-success) / <alpha-value>)",
        dim: "rgb(var(--theme-success-dim) / <alpha-value>)",
      },
    },
    extend: {
      animation: {
        "ping-1": "ping 0.5s cubic-bezier(0, 0, 0.2, 1) 1",
        "spin-slow": "spin 2s linear infinite;",
        "spin-1/4": "spin-1/4 0.5s cubic-bezier(0, 0, 0.2, 1) 1;",
      },
      keyframes: {
        "spin-1/4": {
          from: { transform: "rotate(0deg)" },
          to: { transform: "rotate(90deg)" },
        },
      },
    },
  },
  plugins: [],
};
export default config;
