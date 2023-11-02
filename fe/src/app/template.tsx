"use client";

import { ThemeContext, useTheme } from "@/contexts/themeContext";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import "./globals.css";
import {Menu} from "@/components/menu";

const RootTemplate = ({ children }: { children: React.ReactNode }) => {
  const themeContextValue = useTheme();
  const queryClient = new QueryClient();
  return (
    <ThemeContext.Provider value={themeContextValue}>
      <QueryClientProvider client={queryClient}>
        <Menu />
        {children}
        <ReactQueryDevtools />
      </QueryClientProvider>
    </ThemeContext.Provider>
  );
};

export default RootTemplate;
