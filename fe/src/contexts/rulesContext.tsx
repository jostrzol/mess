import {
  createContext,
  ReactNode,
  useContext,
  useEffect,
  useState,
} from "react";

export interface RulesIndexEntry {
  filename: string;
  modifiedTimestamp: number;
}

export type RulesIndex = RulesIndexEntry[];

export const RuleFilesContext = createContext<RuleFilesContextValue>(null!);
export type RuleFilesContextValue = {
  index: RulesIndex;
  getContent: (filename: string) => string | null
  saveFile: (filename: string, content: string) => void
  removeFile: (filename: string) => void
};

export const useRuleFiles = () => {
  return useContext(RuleFilesContext);
};

export interface RulesFile {
  filename: string;
  content: string;
  saveAs: (filename: string, newContent: string) => void
}

export const RuleFilesProvider = ({ children }: { children?: ReactNode }) => {
  const [index, setIndex] = useState<RulesIndex>([]);

  useEffect(() => {
    const indexStr = localStorage.getItem("rules-index") || "[]";
    const newIndex: RulesIndex = JSON.parse(indexStr);
    setIndex(newIndex);
  }, []);

  const getContent = (filename: string): string | null => {
    return localStorage.getItem(`rules-content:${filename}`);
  }

  const saveFile = (filename: string, content: string) => {
    localStorage.setItem(`rules-content:${filename}`, content);
    const newIndexEntry: RulesIndexEntry = {
      filename,
      modifiedTimestamp: Date.now(),
    }
    const newIndex = [newIndexEntry, ...index.filter(i => i.filename == filename)]
    setIndex(newIndex)
  }

  const removeFile = (filename: string) => {
    localStorage.removeItem(`rules-content:${filename}`);
    const newIndex = index.filter(i => i.filename == filename)
    setIndex(newIndex)
  }

  return (
    <RuleFilesContext.Provider value={{index, getContent, saveFile, removeFile}}>{children}</RuleFilesContext.Provider>
  );
};
