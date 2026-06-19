import type { HistoryOperationResponse as HistoryOperationDTO } from "../types/dto";
import { UserBase } from "../../users/models/user";
import { IDate } from "../../../shared/types/idate";

export class HistoryOperation {
  id: string;
  createdBy: UserBase;
  createdAt: IDate;
  operationType: number;

  constructor(data?: HistoryOperationDTO) {
    ((this.id = data?.id ?? ""),
      (this.createdBy = new UserBase(data?.createdBy)));
    this.createdAt = new IDate(data?.createdAt ?? null);
    this.operationType = data?.operationType ?? 0;
  }

  toDTO(): HistoryOperationDTO {
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
      case 20:
        return "task created";
      case 21:
        return "task updated";
      case 22:
        return "task deleted";
      case 23:
        return "task note added";
      case 24:
        return "task note updated";
      case 25:
        return "task note deleted";
      case 26:
        return "task attachment added";
      case 27:
        return "task attachment deleted";
      case 28:
        return "task time entry added";
      case 29:
        return "task time entry updated";
      case 30:
        return "task time entry deleted";
      default:
        return "unknown operation";
    }
  }
}
