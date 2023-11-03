"use client";

import { createRoom } from "@/api/room";
import { useQuery } from "@tanstack/react-query";
import { redirect } from "next/navigation";

const CreateRoomPage = () => {
  const { data: room, isSuccess } = useQuery({
    queryKey: ["createRoom"],
    queryFn: createRoom,
  });
  if (isSuccess) {
    redirect(`/rooms/${room.id}`);
  }
  return <h1>Creating new room</h1>;
};

export default CreateRoomPage;
