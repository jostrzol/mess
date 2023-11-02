import { ReactNode } from "react";

type MainProps = {
  children: ReactNode;
};

export const Main = ({ children }: MainProps) => {
  return (
    <main className="h-screen flex justify-center items-center">
      {children}
    </main>
  );
};
