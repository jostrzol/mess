import { useOptions } from "@/contexts/optionContext";

export const Options = () => {
  const {current} = useOptions();
  return <h1 className="h-fit">
    options
  </h1>;
};
