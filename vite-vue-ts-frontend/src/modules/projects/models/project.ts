import type { ProjectResponse as ProjectDTO } from "../types/dto";
import { ProjectType } from "../../project-types/models/project-type";
import { ProjectPriority } from "../../project-priorities/models/project-priority";
import { ProjectStatus } from "../../project-statuses/models/project-status";
import { UserBase } from "../../users/models/user";
import { IDate } from "../../../shared/types/idate";
import type { AllowedProjectOperations } from "../../../shared/types/dto/allowed-project-operations";

export class Project {
  id: string;
  slug: string;
  summary: string;
  description: string;
  type: ProjectType;
  priority: ProjectPriority;
  status: ProjectStatus;
  createdAt: IDate;
  updatedAt: IDate;
  startedAt: IDate;
  finishedAt: IDate;
  dueAt: IDate;
  archivedAt: IDate;
  createdBy: UserBase;
  tasksCount: number;
  permissionsCount: number;
  attachmentsCount: number;
  notesCount: number;
  pagesCount: number;
  historyOperationsCount: number;
  allowedOperations: AllowedProjectOperations;

  constructor(data?: ProjectDTO) {
    this.id = data?.id ?? "";
    this.slug = data?.slug ?? "";
    this.summary = data?.summary ?? "";
    this.description = data?.description ?? "";
    this.type = new ProjectType(data?.type);
    this.priority = new ProjectPriority(data?.priority);
    this.status = new ProjectStatus(data?.status);
    this.createdAt = new IDate(data?.createdAt ?? Date.now());
    this.updatedAt = new IDate(data?.updatedAt ?? null);
    this.startedAt = new IDate(data?.startedAt ?? null);
    this.finishedAt = new IDate(data?.finishedAt ?? null);
    this.dueAt = new IDate(data?.dueAt ?? null);
    this.archivedAt = new IDate(data?.archivedAt ?? null);
    this.createdBy = new UserBase(data?.createdBy);
    this.tasksCount = data?.tasksCount ?? 0;
    this.permissionsCount = data?.permissionsCount ?? 0;
    this.attachmentsCount = data?.attachmentsCount ?? 0;
    this.notesCount = data?.notesCount ?? 0;
    this.pagesCount = data?.pagesCount ?? 0;
    this.historyOperationsCount = data?.historyOperationsCount ?? 0;
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

  toDTO(): ProjectDTO {
    return {
      id: this.id,
      slug: this.slug,
      summary: this.summary,
      description: this.description,
      type: this.type.toDTO(),
      priority: this.priority.toDTO(),
      status: this.status.toDTO(),
      createdAt: this.createdAt.msTimestamp ?? 0,
      createdBy: this.createdBy.toDTO(),
      tasksCount: this.tasksCount,
      permissionsCount: this.permissionsCount,
      attachmentsCount: this.attachmentsCount,
      notesCount: this.notesCount,
      pagesCount: this.pagesCount,
      historyOperationsCount: this.historyOperationsCount,
      allowedOperations: this.allowedOperations,
    };
  }
}

export const MAX_SLUG_LENGTH = 8;
export const MAX_SUMMARY_LENGTH = 128;
