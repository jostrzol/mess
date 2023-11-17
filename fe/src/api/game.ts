import { GameState } from "@/model/game/gameState";
import { StaticData } from "@/model/game/gameStaticData";
import { OptionNode, Route } from "@/model/game/options";
import { Resolution } from "@/model/game/resolution";
import { UUID } from "crypto";
import { MessApi } from "./messApi";
import { GameStateDto, gameStateToModel } from "./schema/game";
import { OptionNodeDto, optionNodeToModel, routeToDto } from "./schema/options";
import { ResolutionDto, resolutionToModel } from "./schema/resoultion";
import { StaticDataDto, staticDataToModel } from "./schema/staticData";
import { throwIfError } from "./utils";

export class GameApi extends MessApi {
  public getStaticData = async (roomId: UUID): Promise<StaticData> => {
    const url_ = this.url("rooms/:id/game/static", { params: { id: roomId } });
    const res = await fetch(url_, { method: "GET", credentials: "include" });
    await throwIfError(res);

    const obj: StaticDataDto = await res.json();
    return staticDataToModel(obj);
  };

  public getState = async (roomId: UUID): Promise<GameState> => {
    const url_ = this.url("rooms/:id/game/state", { params: { id: roomId } });
    const res = await fetch(url_, { method: "GET", credentials: "include" });
    await throwIfError(res);

    const obj: GameStateDto = await res.json();
    return gameStateToModel(obj);
  };

  public getTurnOptions = async (roomId: UUID): Promise<OptionNode | null> => {
    const url_ = this.url("rooms/:id/game/options", { params: { id: roomId } });
    const res = await fetch(url_, { method: "GET", credentials: "include" });
    await throwIfError(res);

    const obj: OptionNodeDto | null = await res.json();
    return obj ? optionNodeToModel(obj) : null;
  };

  public playTurn = async (
    roomId: UUID,
    turn: number,
    route: Route,
  ): Promise<GameState> => {
    const url_ = this.url("rooms/:id/game/turns/:turn", {
      params: { id: roomId, turn: turn },
    });
    const res = await fetch(url_, {
      method: "PUT",
      credentials: "include",
      body: JSON.stringify(routeToDto(route)),
    });
    await throwIfError(res);

    const obj: GameStateDto = await res.json();
    return gameStateToModel(obj);
  };

  public getResolution = async (roomId: UUID): Promise<Resolution> => {
    const url_ = this.url("rooms/:id/game/resolution", {
      params: { id: roomId },
    });
    const res = await fetch(url_, { method: "GET", credentials: "include" });
    await throwIfError(res);

    const obj: ResolutionDto = await res.json();
    return resolutionToModel(obj);
  };
}
