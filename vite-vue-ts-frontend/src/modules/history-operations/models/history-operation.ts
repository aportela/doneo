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
      case 100:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.projectCreated";
      case 101:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.projectUpdated";
      case 102:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.projectDeleted";
      case 110:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.projectNoteAdded";
      case 111:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.projectNoteUpdated";
      case 112:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.projectNoteDeleted";
      case 120:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.projectAttachmentAdded";
      case 122:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.projectAttachmentDeleted";
      case 130:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.projectPermissionAdded";
      case 132:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.projectPermissionDeleted";
      case 200:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.taskCreated";
      case 201:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.taskUpdated";
      case 202:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.taskDeleted";
      case 210:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.taskNoteAdded";
      case 211:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.taskNoteUpdated";
      case 212:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.taskNoteDeleted";
      case 220:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.taskAttachmentAdded";
      case 222:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.taskAttachmentDeleted";
      case 230:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.taskTimeEntryAdded";
      case 231:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.taskTimeEntryUpdated";
      case 232:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.taskTimeEntryDeleted";
      default:
        return "modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.body.columns.operationType.unknownOperationType";
    }
  }
}
