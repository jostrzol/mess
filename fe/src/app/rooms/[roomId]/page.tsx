"use client";

import { RoomApi } from "@/api/room";
import { RoomChanged } from "@/api/schema/event";
import { Button } from "@/components/form/button";
import { useMessApi } from "@/contexts/messApiContext";
import { useRoomWebsocket } from "@/contexts/roomWsContext";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { MdContentCopy } from "react-icons/md";
import { RoomPageParams } from "./layout";
import {ConnectionStatus} from "@/components/connectionStatus";

const RoomPage = ({ params }: RoomPageParams) => {
  const router = useRouter();
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
    router.replace(`/rooms/${room.id}/game`);
  }

  return (
    <>
      <ConnectionStatus />
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
          <span className="opacity-0">
            <MdContentCopy />
          </span>
          <h1>Room</h1>
          <MdContentCopy />
        </div>
        <div className="flex justify-between">
          <p>Players</p>
          <p>{`${room.players}/${room.playersNeeded}`}</p>
        </div>
        <div className="flex justify-between">
          <p>Rules</p>
          <pre>{room.rulesFilename}</pre>
        </div>
        <Button type="button" onClick={() => router.push(`/rooms/${room.id}/rules`)}>
          Edit rules
        </Button>
        <Button disabled={!room.isStartable} type="submit">
          Start
        </Button>
      </form>
    </>
  );
};

export default RoomPage;
