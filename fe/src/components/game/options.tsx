import { useOptions } from "@/contexts/optionContext";
import clsx from "clsx";

export const Options = () => {
  const { current } = useOptions();
  return (
    <div className="relative h-10">
      <div className="grid grid-flow-col auto-cols-fr">
        {[...current].map(
          (option, i) => (
            <Option key={i} message={option.message} />
          ),
        )}
      </div>
    </div>
  );
};

const Option = ({ message }: { message: string }) => (
  <div
    className={clsx(
      "mx-2",
      "[&:first-child_hr.left]:invisible",
      "[&:last-child_hr.right]:invisible",
      "w-full",
      "flex",
      "flex-col",
      "align-center",
      "text-sm",
      "text-center",
    )}
  >
    <div className="relative w-full">
      <div className="w-5 h-5 mx-auto border-2 border-primary rounded-full bg-background" />
      <hr className="left absolute w-1/2 top-1/2 left-0 -z-10 -translate-y-1/2 border-b-2" />
      <hr className="right absolute w-1/2 top-1/2 right-0 -z-10 -translate-y-1/2 border-b-2" />
      <div className="absolute w-full h-full bg-background -z-10" />
    </div>
    {message}
  </div>
);
