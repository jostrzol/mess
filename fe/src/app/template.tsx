"use client";

import { Main } from "@/components/main";
import { Menu } from "@/components/menu";
import { ThemeContext, useTheme } from "@/contexts/themeContext";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import "./globals.css";

const RootTemplate = ({ children }: { children: React.ReactNode }) => {
  const themeContextValue = useTheme();
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 60_000,
        throwOnError: true,
      },
      mutations: {
        throwOnError: true,
      },
    },
  });
  return (
    <ThemeContext.Provider value={themeContextValue}>
      <QueryClientProvider client={queryClient}>
        <Menu />
        <Main>{children}</Main>
        <ReactQueryDevtools />
      </QueryClientProvider>
    </ThemeContext.Provider>
  );
};

export default RootTemplate;
