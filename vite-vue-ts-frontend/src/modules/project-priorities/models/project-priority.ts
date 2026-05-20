import type { ProjectPriorityResponse as ProjectPriorityDTO } from "../types/dto";

export class ProjectPriority {
  id: string;
  name: string;
  hexColor: string;

  constructor(data: ProjectPriorityDTO) {
    this.id = data.id;
    this.name = data.name;
    this.hexColor = data.hexColor;
  }

  toDTO(): ProjectPriorityDTO {
    return {
      id: this.id,
      name: this.name,
      hexColor: this.hexColor,
    };
  }
}

export const maxNameLength = 32;
