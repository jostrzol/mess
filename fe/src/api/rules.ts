import { MessApi } from "./messApi";

export class RulesApi extends MessApi {
  public format = async (src: string): Promise<string> => {
    const res = await fetch("/rules/format", {
      method: "PUT",
      credentials: "include",
      body: src,
    });

    return await res.text();
  };
}
