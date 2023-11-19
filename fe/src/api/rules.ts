import { MessApi } from "./messApi";
import { throwIfError } from "./utils";

export class RulesApi extends MessApi {
  public format = async (src: string): Promise<string> => {
    const url_ = this.url("/rules/format");
    const res = await fetch(url_, { method: "PUT", credentials: "include", body: src });
    await throwIfError(res);

    return await res.text();
  };
}
