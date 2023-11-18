type Query = { [key: string]: any };
type Params = Query;
type Options = {
  params?: Params;
  query?: Query;
  schema?: "http" | "ws";
};

export class MessApi {
  _baseUrl: string;

  constructor(baseUrl: string) {
    this._baseUrl = baseUrl;
  }

  public url = (path: string, options?: Options): string => {
    const { params = {}, query = {}, schema = "http" } = options ?? {};

    let url = new URL(`${schema}://${this._baseUrl}`);

    const pathParts = path.split("/").filter((s, i) => i !== 0 || s.length !== 0);
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
    url.pathname += injectedPathParts.join("/");

    Object.entries(query).forEach(([key, value]) => {
      url.searchParams.set(key, value);
    });

    return url.href;
  };
}
