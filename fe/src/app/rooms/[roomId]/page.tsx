"use client";

import { joinRoom } from "@/api/room";
import { Button } from "@/components/form/button";
import { Loader } from "@/components/loader";
import { Main } from "@/components/main";
import { useQuery } from "@tanstack/react-query";
import { UUID } from "crypto";
import { redirect } from "next/navigation";

type RoomPageParams = {
  params: {
    roomId: UUID;
  };
};

const RoomPage = ({ params }: RoomPageParams) => {
  const { data: room, isSuccess } = useQuery({
    queryKey: ["room", params.roomId],
    queryFn: () => joinRoom(params.roomId),
  });
  if (!isSuccess) {
    return <Loader />;
  }
  return (
    <Main>
      <form
        className="w-60 flex flex-col items-stretch gap-4"
        action={async () => {
          redirect(`/rooms/${room.id}/game`);
        }}
      >
        <h1 className="text-center">Room</h1>
        <div className="flex justify-between">
          <p>Players</p>
          <p>{`${room.players}/2`}</p>
        </div>
        <Button disabled={!room.isReady()} type="submit">
          Start
        </Button>
      </form>
    </Main>
  );
};

export default RoomPage;
