"use client";

import { MessApi } from "@/api/messApi";
import { ReactNode, createContext, useContext } from "react";

export const MessApiContext = createContext<string>(null!);

export const useMessApi = <T extends MessApi>(type: {
  new (baseUrl: string): T;
}): T => {
  const messApi = useContext(MessApiContext);
  return new type(messApi);
};

export const MessApiProvider = ({
  baseUrl,
  children,
}: {
  baseUrl: string;
  children?: ReactNode;
}) => {
  return (
    <MessApiContext.Provider value={baseUrl}>
      {children}
    </MessApiContext.Provider>
  );
};
