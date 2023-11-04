const address = "localhost:4000";

type Query = { [key: string]: any };
type Params = Query;
type Options = {
  params?: Params;
  query?: Query;
  schema?: "http" | "ws";
};

export const url = (path: string, options?: Options): string => {
  const {params = {}, query = {}, schema = "http"} = options ?? {}

  let url = new URL(`${schema}://${address}`);

  const pathParts = path.split("/");
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

  Object.entries(query).forEach(([key, value]) => {
    url.searchParams.set(key, value);
  });

  return url.href;
};
