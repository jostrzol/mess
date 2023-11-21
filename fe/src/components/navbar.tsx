import {ReactNode, useState} from "react";
import {IconButton} from "./iconButton";
import {MdArrowBack, MdMenu} from "react-icons/md";
import {useRouter} from "next/navigation";
import {ConnectionStatus} from "./connectionStatus";

export const Navbar = ({children}: {children?: ReactNode}) => {
  const [isMenuOpen, setIsMenuOpen] = useState<boolean>(false)
  const router = useRouter()
  return (
    <div className="w-full p-2 flex gap-2 items-center">
      <IconButton onClick={router.back}><MdArrowBack /></IconButton>
      <IconButton ><MdMenu /></IconButton>
      <div className="grow">
        {children}
      </div>
      <ConnectionStatus className="pr-2" />
    </div>
  )
}

export const NavbarSpacer = () => <div className="h-14 w-full shrink"/>
