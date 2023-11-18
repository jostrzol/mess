import { MessApi } from "@/api/messApi";
import { StaticData } from "@/model/game/gameStaticData";
import { ReactNode, createContext, useContext } from "react";
import { useMessApi } from "./messApiContext";
import {UUID} from "crypto";

export const StaticDataContext = createContext<StaticDataContextValue>(null!);

export const useStaticData = () => {
  return useContext(StaticDataContext);
};

export interface StaticDataContextValue extends StaticData {
  assetUrl: (assetKey: string) => string;
}

export const StaticDataProvider = ({
  roomId,
  staticData,
  children,
}: {
  roomId: UUID;
  staticData: StaticData;
  children?: ReactNode;
}) => {
  const messApi = useMessApi(MessApi);
  const assetUrl = (assetKey: string): string => {
    const url = messApi.url("/rooms/:id/game/assets" + assetKey, {
      params: { id: roomId },
    });
    return url
  };
  return (
    <StaticDataContext.Provider value={{ assetUrl, ...staticData }}>
      {children}
    </StaticDataContext.Provider>
  );
};
