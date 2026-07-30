import type { PagerQuery, PagerResult } from "../../../shared/types/pager";
import type { Order } from "../../../shared/types/order";

import type { ProjectTypeResponse } from "../../project-types/types/dto";
import type { ProjectPriorityResponse } from "../../project-priorities/types/dto";
import type { ProjectStatusResponse } from "../../project-statuses/types/dto";
import type { UserBaseResponse } from "../../users/types/dto";
import type { TimestampRange } from "../../../shared/composables/timestamps";
import type { AllowedProjectOperations } from "../../../shared/types/dto/allowed-project-operations";

export type AddRequest = {
  slug: string;
  summary: string;
  description: string | null;
  type: {
    id: string;
  };
  priority: {
    id: string;
  };
  status: {
    id: string;
  };
};

export type UpdateRequest = {
  id: string;
  slug: string;
  summary: string;
  description: string | null;
  type: {
    id: string;
  };
  priority: {
    id: string;
  };
  status: {
    id: string;
  };
  startedAt: number | null;
  finishedAt: number | null;
  dueAt: number | null;
};

export type PatchRequest = {
  id: string;
  status: {
    id: string;
  };
};

type SearchRequestFilter = {
  slug?: string;
  summary?: string;
  typeId?: string;
  priorityId?: string;
  statusId?: string;
  createdByUserId?: string;
  createdAt?: TimestampRange;
  updatedAt?: TimestampRange;
  finishedAt?: TimestampRange;
  dueAt?: TimestampRange;
  archivedAt?: TimestampRange;
};

export type SearchRequest = {
  pager: PagerQuery;
  order: Order;
  filter?: SearchRequestFilter;
};

export type ProjectResponse = {
  id: string;
  slug: string;
  summary: string;
  description: string;
  type: ProjectTypeResponse;
  priority: ProjectPriorityResponse;
  status: ProjectStatusResponse;
  createdAt: number;
  createdBy: UserBaseResponse;
  updatedAt?: number;
  startedAt?: number;
  finishedAt?: number;
  dueAt?: number;
  archivedAt?: number;
  tasksCount: number;
  permissionsCount: number;
  attachmentsCount: number;
  notesCount: number;
  pagesCount: number;
  historyOperationsCount: number;
  allowedOperations: AllowedProjectOperations;
};

export type SearchResponse = {
  projects: ProjectResponse[];
  pager: PagerResult;
};
