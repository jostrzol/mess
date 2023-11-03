"use client";

import { Button } from "@/components/form/button";
import { Input } from "@/components/form/input";
import { Logo } from "@/components/logo";
import clsx from "clsx";
import { useRouter } from "next/navigation";
import { useState } from "react";

const RootPage = () => {
  const [roomId, setRoomId] = useState("");
  const router = useRouter();
  return (
    <div
      className={clsx(
        "h-full",
        "flex",
        "max-w-lg",
        "flex-col",
        "items-stretch",
        "gap-4",
      )}
    >
      <div
        className={clsx("mx-auto", "m-4", "flex-grow", "flex", "items-center")}
      >
        <Logo size={180} />
      </div>
      <Button onClick={() => router.push("/rooms/")}>New room</Button>
      <form
        className={clsx("flex", "gap-4")}
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
  );
};

export default RootPage;
