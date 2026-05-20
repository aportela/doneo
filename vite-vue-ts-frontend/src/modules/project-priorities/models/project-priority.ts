import type { ProjectPriorityResponse as ProjectPriorityDTO } from "../types/dto";

export class ProjectPriority {
  id: string;
  name: string;
  hexColor: string;
  index: number;

  constructor(data: ProjectPriorityDTO) {
    this.id = data.id;
    this.name = data.name;
    this.hexColor = data.hexColor;
    this.index = data.index;
  }

  toDTO(): ProjectPriorityDTO {
    return {
      id: this.id,
      name: this.name,
      hexColor: this.hexColor,
      index: this.index,
    };
  }
}

export const maxNameLength = 32;
