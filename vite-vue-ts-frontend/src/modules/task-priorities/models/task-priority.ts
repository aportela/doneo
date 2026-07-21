import type { TaskPriorityResponse as TaskPriorityDTO } from "../types/dto";

export class TaskPriority {
  id: string;
  name: string;
  hexColor: string;
  index: number;

  constructor(data?: TaskPriorityDTO) {
    this.id = data?.id ?? "";
    this.name = data?.name ?? "";
    this.hexColor = data?.hexColor ?? "";
    this.index = data?.index ?? 0;
  }

  toDTO(): TaskPriorityDTO {
    return {
      id: this.id,
      name: this.name,
      hexColor: this.hexColor,
      index: this.index,
    };
  }
}

export const MAX_NAME_LENGTH = 32;
