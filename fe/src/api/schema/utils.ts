export const renameKeys = (obj: any): any => {
  if (Array.isArray(obj)) {
    return obj.map((o) => renameKeys(o));
  } else if (typeof obj === "object" && obj !== null) {
    return Object.entries(obj).reduce(
      (r, [k, v]) => ({ ...r, [pascalToSnake(k)]: renameKeys(v) }),
      {}
    );
  } else {
    return obj;
  }
}

const pascalToSnake = (str: string): string => {
  const firstLetter = str.substring(0, 1)
  const rest = str.substring(1)
  return firstLetter.toLowerCase().concat(rest)
}
