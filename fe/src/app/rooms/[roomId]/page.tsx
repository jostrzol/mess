"use client";

import { startGame } from "@/api/game";
import { joinRoom } from "@/api/room";
import { Button } from "@/components/form/button";
import { RoomWsContext } from "@/contexts/roomWsContext";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { redirect } from "next/navigation";
import { useContext, useEffect } from "react";
import { RoomPageParams } from "./layout";

const RoomPage = ({ params }: RoomPageParams) => {
  const client = useQueryClient();
  const { data: room, isSuccess } = useQuery({
    queryKey: ["room", params.roomId],
    queryFn: () => joinRoom(params.roomId),
  });
  const { mutate } = useMutation({
    mutationKey: ["room", params.roomId],
    mutationFn: () => startGame(params.roomId),
    onSuccess: (room) => {
      client.setQueryData(["room", params.roomId], room);
    },
  });

  const { lastEvent } = useContext(RoomWsContext);
  useEffect(() => {
    if (
      lastEvent?.EventType === "RoomChanged" ||
      lastEvent?.EventType === "GameStarted"
    ) {
      client.invalidateQueries({ queryKey: ["room", params.roomId] });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [lastEvent]);

  if (!isSuccess) {
    return null;
  }
  if (room.isStarted) {
    redirect(`/rooms/${room.id}/game`);
  }

  return (
    <>
      <form
        className="w-60 flex flex-col items-stretch gap-4"
        action={() => mutate()}
      >
        <h1 className="text-center">Room</h1>
        <div className="flex justify-between">
          <p>Players</p>
          <p>{`${room.players}/${room.playersNeeded}`}</p>
        </div>
        <Button disabled={!room.isStartable} type="submit">
          Start
        </Button>
      </form>
    </>
  );
};

export default RoomPage;
