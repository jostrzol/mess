"use client";

import { ThemeContext, useTheme } from "@/contexts/themeContext";
import "./globals.css";


const RootTemplate = ({ children }: { children: React.ReactNode }) => {
  const themeContextValue = useTheme();
  return (
      <ThemeContext.Provider value={themeContextValue}>
        {children}
      </ThemeContext.Provider>
  );
};

export default RootTemplate;
