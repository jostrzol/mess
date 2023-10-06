import type {Config} from 'tailwindcss';

const config: Config = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    colors: {
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
        black: "rgb(var(--theme-player-black) / <alpha-value>)",
      },
    },
    extend: {},
  },
  plugins: [],
};
export default config;
