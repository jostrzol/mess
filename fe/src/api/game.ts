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

export class GameApi extends MessApi {
  public getStaticData = async (roomId: UUID): Promise<StaticData> => {
    const res = await this.fetch("rooms/:id/game/static", {
      method: "GET",
      params: { id: roomId },
      credentials: "include",
    });

    const obj: StaticDataDto = await res.json();
    return staticDataToModel(obj);
  };

  public getState = async (roomId: UUID): Promise<GameState> => {
    const res = await this.fetch("rooms/:id/game/state", {
      method: "GET",
      params: { id: roomId },
      credentials: "include",
    });

    const obj: GameStateDto = await res.json();
    return gameStateToModel(obj);
  };

  public getTurnOptions = async (roomId: UUID): Promise<OptionNode | null> => {
    const res = await this.fetch("rooms/:id/game/options", {
      method: "GET",
      params: { id: roomId },
      credentials: "include",
    });

    const obj: OptionNodeDto | null = await res.json();
    return obj ? optionNodeToModel(obj) : null;
  };

  public playTurn = async (
    roomId: UUID,
    turn: number,
    route: Route,
  ): Promise<GameState> => {
    const res = await this.fetch("rooms/:id/game/turns/:turn", {
      method: "PUT",
      params: { id: roomId, turn },
      credentials: "include",
      body: JSON.stringify(routeToDto(route)),
    });

    const obj: GameStateDto = await res.json();
    return gameStateToModel(obj);
  };

  public getResolution = async (roomId: UUID): Promise<Resolution> => {
    const res = await this.fetch("rooms/:id/game/resolution", {
      method: "GET",
      params: { id: roomId },
      credentials: "include",
    });

    const obj: ResolutionDto = await res.json();
    return resolutionToModel(obj);
  };
}
