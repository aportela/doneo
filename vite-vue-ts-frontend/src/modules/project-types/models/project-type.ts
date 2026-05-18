import type { ProjectTypeResponse as ProjectTypeDTO } from "../types/dto";

export class ProjectType {
  id: string;
  name: string;
  hexColor: string;

  constructor(data: ProjectTypeDTO) {
    this.id = data.id;
    this.name = data.name;
    this.hexColor = data.hexColor;
  }

  toDTO(): ProjectTypeDTO {
    return {
      id: this.id,
      name: this.name,
      hexColor: this.hexColor,
    };
  }
}

export const maxNameLength = 32;
