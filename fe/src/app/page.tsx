"use client";

import { Button } from "@/components/form/button";
import { Input } from "@/components/form/input";
import { Main } from "@/components/main";
import clsx from "clsx";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useState } from "react";

const RootPage = () => {
  const [roomId, setRoomId] = useState("");
  const router = useRouter();
  return (
    <Main>
      <div
        className={clsx(
          "h-full",
          "flex",
          "max-w-lg",
          "flex-col",
          "items-stretch",
        )}
      >
        <div className={clsx("mx-auto", "m-4", "flex-grow", "flex", "items-center")}>
          <Image
            src="./favicon.svg"
            alt="logo"
            width={180}
            height={180}
            priority={true}
            className={clsx("h-fit")}
          />
        </div>
        <Button onClick={() => router.push("/rooms/")}>New room</Button>
        <form
          className={clsx("flex")}
          onSubmit={(e) => {
            e.preventDefault();
            router.push(`/rooms/${roomId}`);
          }}
        >
          <Input
            onChange={(e) => setRoomId(e.target.value)}
            placeholder="Room id"
          />
          <Button disabled={roomId === ""} type="submit">
            Join room
          </Button>
        </form>
        <div className={clsx("flex-grow", "h-36", "m-4")} />
      </div>
    </Main>
  );
};

export default RootPage;
