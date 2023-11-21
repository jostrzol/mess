import clsx from "clsx";
import { ReactNode } from "react";

type MainProps = {
  children: ReactNode;
  className?: string;
};

export const Main = ({ children, className }: MainProps) => {
  return (
    <main
      className={clsx(
        "grow relative overflow-hidden",
        "flex flex-col justify-center items-center gap-4",
        className,
      )}
    >
      {children}
    </main>
  );
};
