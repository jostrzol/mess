import { StaticData } from "@/model/game/gameStaticData";
import { ReactNode, createContext, useContext } from "react";

export const StaticDataContext = createContext<StaticData>(null!);

export const useStaticData = () => {
  return useContext(StaticDataContext);
};

export const StaticDataProvider = ({
  staticData,
  children,
}: {
  staticData: StaticData;
  children?: ReactNode;
}) => {
  return (
    <StaticDataContext.Provider value={staticData}>
      {children}
    </StaticDataContext.Provider>
  );
};
