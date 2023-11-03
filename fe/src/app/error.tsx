"use client"; // Error components must be Client Components

import { Button } from "@/components/form/button";
import { Main } from "@/components/main";
import { useEffect } from "react";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <Main>
      <div className="flex flex-col items-center gap-4">
        <h1>Something went wrong!</h1>
        <Button
          onClick={
            // Attempt to recover by trying to re-render the segment
            () => reset()
          }
        >
          Try again
        </Button>
      </div>
    </Main>
  );
}
