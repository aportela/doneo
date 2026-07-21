import type { TaskStatusResponse as TaskStatusDTO } from "../types/dto";
import type { StatusFlags } from "../../../shared/types/status-flags";

export class TaskStatus {
  id: string;
  name: string;
  hexColor: string;
  index: number;
  flags: StatusFlags;

  constructor(data?: TaskStatusDTO) {
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

  toDTO(): TaskStatusDTO {
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
