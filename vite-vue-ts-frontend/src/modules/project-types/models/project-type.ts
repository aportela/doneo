import type {
  AddRequest,
  ProjectTypeResponse as ProjectTypeDTO,
  UpdateRequest,
} from "../types/dto";

export class ProjectType {
  id: string;
  name: string;
  hexColor: string;

  constructor(data?: ProjectTypeDTO) {
    this.id = data?.id ?? "";
    this.name = data?.name ?? "";
    this.hexColor = data?.hexColor ?? "";
  }

  toDTO(): ProjectTypeDTO {
    return {
      id: this.id,
      name: this.name,
      hexColor: this.hexColor,
    };
  }

  toAddProjectTypeRequestPayload(): AddRequest {
    return {
      name: this.name,
      HexColor: this.hexColor,
    };
  }

  toUpdateProjectTypeRequestPayload(): UpdateRequest {
    return {
      id: this.id,
      name: this.name,
      HexColor: this.hexColor,
    };
  }
}

export const MAX_NAME_LENGTH = 32;
