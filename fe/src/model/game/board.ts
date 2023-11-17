import { Square } from "./square";

export interface Board {
  height: number;
  width: number;
}

export namespace Board {
  export const MapSquares = <T>(
    board: Board,
    func: (square: Square, key: string) => T,
  ): T[] =>
    [...Array(board.height).keys()].flatMap((_, j) =>
      [...Array(board.width).keys()].map((_, i) => {
        const y = board.height - 1 - j;
        const x = i;
        const square: [number, number] = [x, y];
        const key = Square.toString(square);
        return func(square, key);
      }),
    );
}
