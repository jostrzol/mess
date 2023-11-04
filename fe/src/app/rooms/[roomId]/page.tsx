"use client";

import { joinRoom, useRoomWebsocket } from "@/api/room";
import { ConnectionStatus } from "@/components/connectionStatus";
import { Button } from "@/components/form/button";
import { Loader } from "@/components/loader";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { UUID } from "crypto";
import { redirect } from "next/navigation";
import { useEffect } from "react";

type RoomPageParams = {
  params: {
    roomId: UUID;
  };
};

const RoomPage = ({ params }: RoomPageParams) => {
  const client = useQueryClient()
  const { data: room, isSuccess } = useQuery({
    queryKey: ["room", params.roomId],
    queryFn: () => joinRoom(params.roomId),
  });

  const { lastEvent, readyState } = useRoomWebsocket(
    params.roomId,
  );

  useEffect(() => {
    if (lastEvent?.EventType === "RoomChanged") {
      client.invalidateQueries({queryKey: ["room", params.roomId]})
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [lastEvent]);

  if (!isSuccess) {
    return <Loader />;
  }
  return (
    <>
      <ConnectionStatus state={readyState} />
      <form
        className="w-60 flex flex-col items-stretch gap-4"
        action={() => {
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
    </>
  );
};

export default RoomPage;
