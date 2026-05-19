import type { ProjectStatusResponse as ProjectStatusDTO } from "../types/dto";

export class ProjectStatus {
  id: string;
  name: string;
  hexColor: string;
  index: number;

  constructor(data: ProjectStatusDTO) {
    this.id = data.id;
    this.name = data.name;
    this.hexColor = data.hexColor;
    this.index = data.index;
  }

  toDTO(): ProjectStatusDTO {
    return {
      id: this.id,
      name: this.name,
      hexColor: this.hexColor,
      index: this.index,
    };
  }
}

export const maxNameLength = 32;
