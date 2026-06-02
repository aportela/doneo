import type { ProjectHistoryOperationResponse as ProjectHistoryOperationDTO } from "../types/dto";
import { UserBase } from "../../users/models/user";
import { IDate } from "../../../shared/types/idate";

export class ProjectHistoryOperation {
  createdBy: UserBase;
  createdAt: IDate;
  operationType: number;

  constructor(data?: ProjectHistoryOperationDTO) {
    this.createdBy = new UserBase(data?.createdBy);
    this.createdAt = new IDate(data?.createdAt ?? null);
    this.operationType = data?.operationType ?? 0;
  }

  toDTO(): ProjectHistoryOperationDTO {
    return {
      createdBy: this.createdBy.toDTO(),
      createdAt: this.createdAt?.msTimestamp ?? 0,
      operationType: this.operationType ?? 0,
    };
  }

  ToString(): string {
    switch (this.operationType) {
      case 1:
        return "project created";
      case 2:
        return "project updated";
      default:
        return "unknown operation";
    }
  }
}
