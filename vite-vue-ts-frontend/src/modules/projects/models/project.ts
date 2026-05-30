import type { ProjectResponse as ProjectDTO } from "../types/dto";
import { ProjectType } from "../../project-types/models/project-type";
import { ProjectPriority } from "../../project-priorities/models/project-priority";
import { ProjectStatus } from "../../project-statuses/models/project-status";
import { UserBase } from "../../users/models/user";
import { IDate } from "../../../shared/types/idate";

export class Project {
  id: string | null;
  key: string | null;
  summary: string | null;
  description: string | null;
  type: ProjectType;
  priority: ProjectPriority;
  status: ProjectStatus;
  createdAt: IDate;
  updatedAt?: IDate;
  startedAt?: IDate;
  finishedAt?: IDate;
  dueAt?: IDate;
  createdBy: UserBase;
  tasksCount: number;
  permissionsCount: number;
  attachmentsCount: number;
  notesCount: number;
  historyOperationsCount: number;

  constructor(data?: ProjectDTO) {
    this.id = data?.id ?? null;
    this.key = data?.key ?? null;
    this.summary = data?.summary ?? null;
    this.description = data?.description ?? null;
    this.type = new ProjectType(data?.type);
    this.priority = new ProjectPriority(data?.priority);
    this.status = new ProjectStatus(data?.status);
    this.createdAt = new IDate(data?.createdAt ?? new Date().getTime());
    this.updatedAt = data?.updatedAt ? new IDate(data?.updatedAt) : undefined;
    this.startedAt = data?.startedAt ? new IDate(data?.startedAt) : undefined;
    this.finishedAt = data?.finishedAt
      ? new IDate(data?.finishedAt)
      : undefined;
    this.dueAt = data?.dueAt ? new IDate(data?.dueAt) : undefined;
    this.createdBy = new UserBase(data?.createdBy);
    this.tasksCount = data?.tasksCount ?? 0;
    this.permissionsCount = data?.permissionsCount ?? 0;
    this.attachmentsCount = data?.attachmentsCount ?? 0;
    this.notesCount = data?.notesCount ?? 0;
    this.historyOperationsCount = data?.historyOperationsCount ?? 0;
  }

  toDTO(): ProjectDTO {
    return {
      id: this.id ?? "",
      key: this.key ?? "",
      summary: this.summary ?? "",
      description: this.description ?? "",
      type: this.type.toDTO(),
      priority: this.priority.toDTO(),
      status: this.status.toDTO(),
      createdAt: this.createdAt.msTimestamp,
      createdBy: this.createdBy.toDTO(),
      tasksCount: this.tasksCount,
      permissionsCount: this.permissionsCount,
      attachmentsCount: this.attachmentsCount,
      notesCount: this.notesCount,
      historyOperationsCount: this.historyOperationsCount,
    };
  }
}

export const MAX_KEY_LENGTH = 8;
export const MAX_SUMMARY_LENGTH = 128;
