import type { TaskStatusResponse as TaskStatusDTO } from "../types/dto";

export class TaskStatus {
  id: string;
  name: string;
  hexColor: string;

  constructor(data: TaskStatusDTO) {
    this.id = data.id;
    this.name = data.name;
    this.hexColor = data.hexColor;
  }

  toDTO(): TaskStatusDTO {
    return {
      id: this.id,
      name: this.name,
      hexColor: this.hexColor,
    };
  }
}

export const MAX_NAME_LENGTH = 32;
