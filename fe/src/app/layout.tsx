import { MessApiProvider } from "@/contexts/messApiContext";
import { Theme, defaultTheme } from "@/model/theme";
import Color from "color";
import type { Metadata } from "next";
import { Inter } from "next/font/google";
import { cookies } from "next/headers";
import config from "../../app.config";
import "./globals.css";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "mess",
  description: "A chess-like custom game engine.",
};

const RootLayout = ({ children }: { children: React.ReactNode }) => {
  const themeCookie = cookies().get("theme");
  const theme: Theme =
    themeCookie !== undefined ? JSON.parse(themeCookie?.value) : defaultTheme;
  const colorDefinitions = Object.entries(theme.colors).map(([key, value]) => {
    const colorKey = `--theme-${key}`;
    const colorValue = Color(value).array().join(" ");
    return `${colorKey}: ${colorValue};`;
  });
  const style = `:root { ${colorDefinitions.join("")} }`;
  return (
    <html lang="en" className={inter.className}>
      <head>
        <style>{style}</style>
      </head>
      <body className="bg-background text-txt flex flex-col h-screen">
        <MessApiProvider baseUrl={config.backendUrl}>
          {children}
        </MessApiProvider>
      </body>
    </html>
  );
};

export default RootLayout;
