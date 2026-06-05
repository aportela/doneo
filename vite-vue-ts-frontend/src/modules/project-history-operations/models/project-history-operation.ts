import type { ProjectHistoryOperationResponse as ProjectHistoryOperationDTO } from "../types/dto";
import { UserBase } from "../../users/models/user";
import { IDate } from "../../../shared/types/idate";

export class ProjectHistoryOperation {
  id: string;
  createdBy: UserBase;
  createdAt: IDate;
  operationType: number;

  constructor(data?: ProjectHistoryOperationDTO) {
    ((this.id = data?.id ?? ""),
      (this.createdBy = new UserBase(data?.createdBy)));
    this.createdAt = new IDate(data?.createdAt ?? null);
    this.operationType = data?.operationType ?? 0;
  }

  toDTO(): ProjectHistoryOperationDTO {
    return {
      id: this.id ?? "",
      createdBy: this.createdBy.toDTO(),
      createdAt: this.createdAt?.msTimestamp ?? 0,
      operationType: this.operationType ?? 0,
    };
  }

  getOperationTypeLabel(): string {
    switch (this.operationType) {
      case 1:
        return "project created";
      case 2:
        return "project updated";
      case 3:
        return "project deleleted";
      case 4:
        return "project note added";
      case 5:
        return "project note updated";
      case 6:
        return "project note deleted";
      case 7:
        return "project attachment added";
      case 8:
        return "project attachment deleted";
      case 9:
        return "project permission added";
      case 10:
        return "project permission deleted";
      default:
        return "unknown operation";
    }
  }
  getNaiveUITimelineItemType():
    | "default"
    | "success"
    | "info"
    | "warning"
    | "error" {
    switch (this.operationType) {
      case 1: // add project
        return "success";
      case 2: // update project
        return "info";
      case 3: // delete project
        return "error";
      case 4: // add note
        return "success";
      case 5: // update note
        return "info";
      case 6: // delete note
        return "error";
      case 7: // added attachment
        return "success";
      case 8: // delete attachment
        return "error";
      case 9: // added permission
        return "success";
      case 10: // delete permission
        return "error";
      default:
        return "default";
    }
  }
}
