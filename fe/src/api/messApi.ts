type Query = { [key: string]: any };
type Params = Query;
type Options = {
  params?: Params;
  query?: Query;
  schema?: "http" | "ws";
};

type FetchOptions = Parameters<typeof fetch>[1];

export class MessApi {
  _baseUrl: string;

  constructor(baseUrl: string) {
    this._baseUrl = baseUrl;
  }

  public fetch = async (
    path: string,
    options: FetchOptions & Options,
  ): ReturnType<typeof fetch> => {
    const { params, query, schema, ...fetchOptions } = options;
    const url = this.url(path, { params, query, schema });
    const res = fetch(url, fetchOptions).then(async (res) => {
      if (!res.ok) {
        const status = `${res.status} ${res.statusText}`;
        const body = await res.text();
        throw new InvalidServerResponseError(status, body);
      } else {
        return res;
      }
    });

    try {
      return await res;
    } catch (e) {
      let error: any;
      if (e instanceof InvalidServerResponseError) {
        error = { response: e.toObject(), request: { url, ...fetchOptions } };
      } else if (e instanceof Error) {
        error = { cause: e.message, request: { url, ...fetchOptions } };
      } else {
        error = { cause: `${e}`, request: { url, ...fetchOptions } };
      }
      const errorStr = JSON.stringify(error, null, 4);
      throw new Error(`failed to fetch: ${errorStr}`, { cause: e });
    }
  };

  public url = (path: string, options?: Options): string => {
    const { params = {}, query = {}, schema = "http" } = options ?? {};

    let url = new URL(`${schema}://${this._baseUrl}`);

    const pathParts = path
      .split("/")
      .filter((s, i) => i !== 0 || s.length !== 0);
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

class InvalidServerResponseError implements Error {
  get name() {
    return "InvalidServerResponseError";
  }
  status: string;
  body: string;
  get message() {
    return `invalid server response: ${JSON.stringify(this.toObject())}`;
  }

  constructor(status: string, body: string) {
    this.status = status;
    this.body = body;
  }

  toObject() {
    return {
      status: this.status,
      body: toJSONOrText(this.body),
    };
  }
}

const toJSONOrText = (raw: string): any | string => {
  try {
    return JSON.parse(raw);
  } catch {
    return raw;
  }
};
