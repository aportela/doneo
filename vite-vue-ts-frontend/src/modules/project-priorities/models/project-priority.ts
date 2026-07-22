import type {
  AddRequest,
  ProjectPriorityResponse as ProjectPriorityDTO,
  UpdateRequest,
} from "../types/dto";

export class ProjectPriority {
  id: string;
  name: string;
  hexColor: string;
  index: number;

  constructor(data?: ProjectPriorityDTO) {
    this.id = data?.id ?? "";
    this.name = data?.name ?? "";
    this.hexColor = data?.hexColor ?? "";
    this.index = data?.index ?? 0;
  }

  toDTO(): ProjectPriorityDTO {
    return {
      id: this.id,
      name: this.name,
      hexColor: this.hexColor,
      index: this.index,
    };
  }

  toAddProjectPriorityRequestPayload(): AddRequest {
    return {
      name: this.name,
      hexColor: this.hexColor,
      index: this.index,
    };
  }

  toUpdateProjectPriorityRequestPayload(): UpdateRequest {
    return {
      id: this.id,
      name: this.name,
      hexColor: this.hexColor,
      index: this.index,
    };
  }
}

export const MAX_NAME_LENGTH = 32;
