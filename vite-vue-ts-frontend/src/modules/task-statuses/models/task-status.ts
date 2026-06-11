import type { TaskStatusResponse as TaskStatusDTO } from "../types/dto";
import type { StatusFlags } from "../../../shared/types/status-flags";

export class TaskStatus {
  id: string | null;
  name: string | null;
  hexColor: string | null;
  index: number;
  flags: StatusFlags;

  constructor(data?: TaskStatusDTO) {
    this.id = data?.id ?? null;
    this.name = data?.name ?? null;
    this.hexColor = data?.hexColor ?? null;
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
      id: this.id ?? "",
      name: this.name ?? "",
      hexColor: this.hexColor ?? "",
      index: this.index ?? 0,
      flags: this.flags ?? {
        defaultStatusOnCreation: false,
        fillEmptyStartDate: false,
        setStartDate: false,
        fillEmptyFinishDate: false,
        setFinishDate: false,
        unsetFinishDateOnLeave: false,
      },
    };
  }
}

export const MAX_NAME_LENGTH = 32;
