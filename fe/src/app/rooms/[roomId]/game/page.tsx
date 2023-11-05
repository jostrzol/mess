"use client";

import { getGame } from "@/api/game";
import { Board } from "@/components/game/board";
import { useQuery } from "@tanstack/react-query";
import { RoomPageParams } from "../layout";

const GamePage = ({ params }: RoomPageParams) => {
  const { data: state, isSuccess } = useQuery({
    queryKey: ["room", params.roomId, "game"],
    queryFn: () => getGame(params.roomId),
  });
  if (!isSuccess) {
    return null
  }
  return (
    <Board
      pieces={state?.pieces}
      board={{
        height: 8,
        width: 8,
      }}
    />
  );
};

export default GamePage;
