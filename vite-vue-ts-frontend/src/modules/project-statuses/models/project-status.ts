import type { ProjectStatusResponse as ProjectStatusDTO } from "../types/dto";
import type { StatusFlags } from "../../../shared/types/status-flags";

export class ProjectStatus {
  id: string;
  name: string;
  hexColor: string;
  index: number;
  flags: StatusFlags;

  constructor(data?: ProjectStatusDTO) {
    this.id = data?.id ?? "";
    this.name = data?.name ?? "";
    this.hexColor = data?.hexColor ?? "";
    this.index = data?.index ?? 0;
    this.flags = data?.flags ?? {
      defaultStatusOnCreation: false,
      fillEmptyStartDate: false,
      setStartDate: false,
      fillEmptyFinishDate: false,
      setFinishDate: false,
      unsetFinishDateOnLeave: false,
    };
  }

  toDTO(): ProjectStatusDTO {
    return {
      id: this.id,
      name: this.name,
      hexColor: this.hexColor,
      index: this.index,
      flags: this.flags,
    };
  }
}

export const MAX_NAME_LENGTH = 32;
