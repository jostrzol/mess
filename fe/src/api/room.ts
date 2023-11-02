import { UUID } from "crypto";

interface RoomDto {
  ID: UUID;
}

export const createRoom = async () => {
  const res = await fetch("http://localhost:4000/rooms", { method: "POST" });

  if (!res.ok) {
    throw new Error("Failed to fetch data");
  }

  const obj: RoomDto = await res.json();
  return {
    id: obj.ID,
  };
}
