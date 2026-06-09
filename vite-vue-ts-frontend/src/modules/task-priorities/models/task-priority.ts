import type { TaskPriorityResponse as TaskPriorityDTO } from "../types/dto";

export class TaskPriority {
  id: string | null;
  name: string | null;
  hexColor: string | null;
  index: number | null;

  constructor(data?: TaskPriorityDTO) {
    this.id = data?.id ?? null;
    this.name = data?.name ?? null;
    this.hexColor = data?.hexColor ?? null;
    this.index = data?.index ?? 0;
  }

  toDTO(): TaskPriorityDTO {
    return {
      id: this.id ?? "",
      name: this.name ?? "",
      hexColor: this.hexColor ?? "",
      index: this.index ?? 0,
    };
  }
}

export const MAX_NAME_LENGTH = 32;
