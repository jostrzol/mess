"use client";

import { RoomChanged } from "@/api/schema/event";
import { Button } from "@/components/form/button";
import { useRoomWebsocket } from "@/contexts/roomWsContext";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { redirect } from "next/navigation";
import { MdContentCopy } from "react-icons/md";
import { RoomPageParams } from "./layout";
import {useMessApi} from "@/contexts/messApiContext";
import {RoomApi} from "@/api/room";

const RoomPage = ({ params }: RoomPageParams) => {
  const roomApi = useMessApi(RoomApi);
  const client = useQueryClient();
  const { data: room, isSuccess } = useQuery({
    queryKey: ["room", params.roomId],
    queryFn: () => roomApi.joinRoom(params.roomId),
  });
  const { mutate } = useMutation({
    mutationKey: ["room", params.roomId],
    mutationFn: () => roomApi.startGame(params.roomId),
    onSuccess: (room) => {
      client.setQueryData(["room", params.roomId], room);
    },
  });

  useRoomWebsocket<RoomChanged>({
    type: "RoomChanged",
    onEvent: () => {
      client.invalidateQueries({ queryKey: ["room", params.roomId] });
    },
  });

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
        <div
          className="flex m-auto gap-2 cursor-pointer"
          onClick={() =>
            navigator.clipboard.writeText(window.location.toString())
          }
        >
          <span className="opacity-0"><MdContentCopy /></span>
          <h1>Room</h1>
          <MdContentCopy />
        </div>
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
