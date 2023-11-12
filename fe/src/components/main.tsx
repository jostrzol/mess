import { ReactNode } from "react";

type MainProps = {
  children: ReactNode;
};

export const Main = ({ children }: MainProps) => {
  return (
    <main className="h-screen p-4 relative overflow-hidden flex flex-col justify-center items-center gap-4">
      {children}
    </main>
  );
};
