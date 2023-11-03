import clsx from "clsx";
import Image from "next/image";
import { useRouter } from "next/navigation";

type LogoProps = {
  size: number;
  className?: string;
};

export const Logo = ({ size, className }: LogoProps) => {
  const router = useRouter();
  return (
    <Image
      src="./favicon.svg"
      alt="logo"
      width={size}
      height={size}
      priority={true}
      className={clsx("rounded-full", "hover:cursor-pointer", className)}
      onClick={() => router.replace("/")}
    />
  );
};
