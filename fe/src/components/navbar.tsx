import { useRouter } from "next/navigation";
import { ReactNode, useState } from "react";
import { MdArrowBack, MdMenu } from "react-icons/md";
import { ConnectionStatus } from "./connectionStatus";
import { IconButton } from "./iconButton";
import { Menu } from "./menu";

export const Navbar = ({ children }: { children?: ReactNode }) => {
  const [isMenuOpen, setIsMenuOpen] = useState<boolean>(false);
  const router = useRouter();
  return (
    <div className="w-full p-2 flex gap-2 items-center">
      <IconButton onClick={router.back}>
        <MdArrowBack />
      </IconButton>
      <IconButton onClick={() => setIsMenuOpen(true)}>
        <MdMenu />
      </IconButton>
      <div className="grow">{children}</div>
      <ConnectionStatus className="pr-2" />
      <Menu open={isMenuOpen} onClose={() => setIsMenuOpen(!isMenuOpen)} />
    </div>
  );
};

export const NavbarSpacer = () => <div className="h-14 w-full shrink" />;
