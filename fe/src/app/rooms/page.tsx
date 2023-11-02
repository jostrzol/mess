import { createRoom } from "@/api/room";
import { redirect } from "next/navigation";

const CreateRoomPage = async () => {
  const room = await createRoom();
  redirect(`/rooms/${room.id}`)
};

export default CreateRoomPage;
