import type { ProjectStatusResponse as ProjectStatusDTO } from "../types/dto";

export class ProjectStatus {
  id: string;
  name: string;
  hexColor: string;

  constructor(data: ProjectStatusDTO) {
    this.id = data.id;
    this.name = data.name;
    this.hexColor = data.hexColor;
  }

  toDTO(): ProjectStatusDTO {
    return {
      id: this.id,
      name: this.name,
      hexColor: this.hexColor,
    };
  }
}

export const MAX_NAME_LENGTH = 32;
