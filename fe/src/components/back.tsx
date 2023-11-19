import clsx from "clsx";
import { useRouter } from "next/navigation";
import { MdArrowBack } from "react-icons/md";

export const Back = () => {
  const router = useRouter();
  return (
    <div
      className={clsx(
        "p-2",
        "self-center",
        "bg-primary-dim/20",
        "hover:bg-primary-dim/40",
        "cursor-pointer",
        "rounded-full",
      )}
      onClick={router.back}
    >
      <MdArrowBack />
    </div>
  );
};
