export type Square = [number, number];

export namespace Square {
  const A_CODE = "A".charCodeAt(0);
  const Z_CODE = "Z".charCodeAt(0);
  const MAX_X = Z_CODE - A_CODE;
  export const file = (square: Square): string => {
    if (square[0] > MAX_X) {
      throw Error(
        `Coord x=${square[0]} too big to stringify. Maximum is ${MAX_X}`,
      );
    }
    const code = A_CODE + square[0];
    return String.fromCharCode(code);
  };

  export const rank = (square: Square): string => {
    return (square[1] + 1).toString();
  };

  export const toString = (square: Square): string =>
    file(square) + rank(square);

  export const equals = (first: Square, second: Square): boolean =>
    first[0] === second[0] && first[1] === second[1];
}
