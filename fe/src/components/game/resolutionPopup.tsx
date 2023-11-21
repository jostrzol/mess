import { Resolution } from "@/model/game/resolution";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { Button } from "../form/button";
import { Popup } from "../popup";

export const ResolutionPopup = ({
  status,
}: {
  status: Resolution["status"];
}) => {
  const [message, icon] = {
    Unresolved: ["???", null],
    Win: ["Victory", "/pieces/king.svg"],
    Draw: ["Draw", "/pieces/knight.svg"],
    Defeat: ["Defeat", "/pieces/pawn.svg"],
  }[status];
  const router = useRouter();
  return (
    <Popup modal>
      <div className="px-4 flex flex-col items-center">
        {icon && (
          <Image alt={"resolution icon"} src={icon} height={32} width={32} />
        )}
        <h1>{message}</h1>
        <section className="mt-8 w-full flex flex-col items-stretch gap-2">
          <Button disabled>Rematch</Button>
          <Button disabled>Return to room</Button>
          <Button onClick={() => router.push("/")}>Exit to main menu</Button>
        </section>
      </div>
    </Popup>
  );
};
