import type {Config} from 'tailwindcss';

const config: Config = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    colors: {
      background: "var(--theme-background)",
      primary: {
        DEFAULT: "var(--theme-primary)",
        dim: "var(--theme-primary-dim)",
      },
      txt: {
        DEFAULT: "var(--theme-txt)",
        dim: "var(--theme-txt-dim)",
      },
      player: {
        white: "var(--theme-player-white)",
        black: "var(--theme-player-black)",
      },
    },
    extend: {},
  },
  plugins: [],
};
export default config;
