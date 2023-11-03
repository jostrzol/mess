const address = "http://localhost:4000";

type Query = { [key: string]: any };
type Params = Query
type Options = {
  params?: Params,
  query?: Query
}

export const url = (path: string, options?: Options): string => {
  let url = new URL(address);

  const pathParts = path.split("/");
  const params = options?.params || {}
  const injectedPathParts = pathParts.map((part) => {
    if (!part.startsWith(":")) {
      return part;
    }
    const value = params[part.slice(1)];
    if (value === undefined) {
      throw new Error(`parameter ${part} not provided`);
    }
    return encodeURIComponent(value);
  });
  url.pathname = injectedPathParts.join("/");

  Object.entries(options?.query || {}).forEach(([key, value]) => {
    url.searchParams.set(key, value);
  });

  return url.href;
};
