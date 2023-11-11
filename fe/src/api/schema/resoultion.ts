import {Resolution} from "@/model/game/resolution"

export interface ResolutionDto {
  Status: "Unresolved" | "Win" | "Draw" | "Defeat"
}

export const resolutionToModel = (resolution: ResolutionDto): Resolution => ({
  status: resolution.Status
})
