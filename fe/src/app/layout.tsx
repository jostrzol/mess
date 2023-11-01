import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "mess",
  description: "A chess-like custom game engine.",
};

const RootLayout = ({ children }: { children: React.ReactNode }) => {
  return (
    <html lang="en" className={inter.className}>
        <body className="bg-background text-txt">{children}</body>
    </html>
  );
};

export default RootLayout;
