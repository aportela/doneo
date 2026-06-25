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

  // TODO: i18n
  getOperationTypeLabel(): string {
    switch (this.operationType) {
      case 100:
        return "project created";
      case 101:
        return "project updated";
      case 102:
        return "project deleted";
      case 110:
        return "project note added";
      case 111:
        return "project note updated";
      case 112:
        return "project note deleted";
      case 120:
        return "project attachment added";
      case 122:
        return "project attachment deleted";
      case 130:
        return "project permission added";
      case 132:
        return "project permission deleted";
      case 200:
        return "task created";
      case 201:
        return "task updated";
      case 202:
        return "task deleted";
      case 210:
        return "task note added";
      case 211:
        return "task note updated";
      case 212:
        return "task note deleted";
      case 220:
        return "task attachment added";
      case 222:
        return "task attachment deleted";
      case 230:
        return "task time entry added";
      case 231:
        return "task time entry updated";
      case 232:
        return "task time entry deleted";
      default:
        return "unknown operation";
    }
  }
}
