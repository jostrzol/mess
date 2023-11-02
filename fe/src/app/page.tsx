"use client";

import { Button } from "@/components/form/button";
import { Input } from "@/components/form/input";
import { Main } from "@/components/main";
import clsx from "clsx";
import { useRouter } from "next/navigation";
import { useState } from "react";

const RootPage = () => {
  const [roomId, setRoomId] = useState("");
  const router = useRouter();
  return (
    <Main>
      <div className={clsx("flex", "max-w-lg", "flex-col")}>
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
      </div>
    </Main>
  );
};

export default RootPage;
