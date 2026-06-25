import type { TaskResponse as ProjectTaskDTO } from "../types/dto";
import { TaskPriority } from "../../task-priorities/models/task-priority";
import { TaskStatus } from "../../task-statuses/models/task-status";
import { UserBase } from "../../users/models/user";
import { IDate } from "../../../shared/types/idate";
import type { AllowedProjectOperations } from "../../../shared/types/dto/allowed-project-operations";

export class Task {
  id: string | null;
  projectId: string | null;
  slug: string | null;
  summary: string | null;
  description: string | null;
  priority: TaskPriority;
  status: TaskStatus;
  createdAt: IDate;
  updatedAt: IDate;
  startedAt: IDate;
  finishedAt: IDate;
  dueAt: IDate;
  totalSpentTime: number;
  estimatedTime: number;
  createdBy: UserBase;
  tags: string[];
  permissionsCount: number;
  attachmentsCount: number;
  notesCount: number;
  historyOperationsCount: number;
  timeTrackingsCount: number;
  allowedOperations: AllowedProjectOperations;

  constructor(data?: ProjectTaskDTO) {
    this.id = data?.id ?? null;
    this.projectId = data?.projectId ?? null;
    this.slug = data?.slug ?? null;
    this.summary = data?.summary ?? null;
    this.description = data?.description ?? null;
    this.priority = new TaskPriority(data?.priority);
    this.status = new TaskStatus(data?.status);
    this.createdAt = new IDate(data?.createdAt ?? new Date().getTime());
    this.updatedAt = new IDate(data?.updatedAt ?? null);
    this.startedAt = new IDate(data?.startedAt ?? null);
    this.finishedAt = new IDate(data?.finishedAt ?? null);
    this.dueAt = new IDate(data?.dueAt ?? null);
    this.totalSpentTime = data?.totalSpentTime ?? 0;
    this.estimatedTime = data?.estimatedTime ?? 0;
    this.createdBy = new UserBase(data?.createdBy);
    this.tags = data?.tags ?? [];
    this.permissionsCount = data?.permissionsCount ?? 0;
    this.attachmentsCount = data?.attachmentsCount ?? 0;
    this.notesCount = data?.notesCount ?? 0;
    this.historyOperationsCount = data?.historyOperationsCount ?? 0;
    this.timeTrackingsCount = data?.timeTrackingsCount ?? 0;
    this.allowedOperations = data?.allowedOperations ?? {
      updateProject: false,
      deleteProject: false,
      viewProject: false,
      addTask: false,
      updateTask: false,
      deleteTask: false,
      viewTask: false,
    };
  }

  toDTO(): ProjectTaskDTO {
    return {
      id: this.id ?? "",
      projectId: this.projectId ?? "",
      slug: this.slug ?? "",
      summary: this.summary ?? "",
      description: this.description ?? "",
      priority: this.priority.toDTO(),
      status: this.status.toDTO(),
      createdAt: this.createdAt.msTimestamp ?? 0,
      createdBy: this.createdBy.toDTO(),
      totalSpentTime: this.totalSpentTime ?? 0,
      estimatedTime: this.estimatedTime ?? 0,
      tags: this.tags,
      permissionsCount: this.permissionsCount,
      attachmentsCount: this.attachmentsCount,
      notesCount: this.notesCount,
      historyOperationsCount: this.historyOperationsCount,
      timeTrackingsCount: this.timeTrackingsCount,
      allowedOperations: this.allowedOperations,
    };
  }
}

export const MAX_SUMMARY_LENGTH = 128;
