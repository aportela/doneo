import type {
  PagerRequest,
  PagerResponse,
} from "../../../shared/types/dto/pager";
import type { Order } from "../../../shared/types/dto/order";

import type { ProjectTypeResponse } from "../../project-types/types/dto";
import type { ProjectPriorityResponse } from "../../project-priorities/types/dto";
import type { ProjectStatusResponse } from "../../project-statuses/types/dto";
import type { UserBaseResponse } from "../../users/types/dto";
import type { TimestampRange } from "../../../shared/composables/timestamps";

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
  createdAt?: TimestampRange;
  createdByUserId?: string;
};

export type SearchRequest = {
  pager: PagerRequest;
  order: Order;
  filter?: SearchRequestFilter;
};

type allowedProjectOperations = {
  viewProject: boolean;
  updateProject: boolean;
  deleteProject: boolean;
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
  tasksCount: number;
  permissionsCount: number;
  attachmentsCount: number;
  notesCount: number;
  historyOperationsCount: number;
  allowedOperations: allowedProjectOperations;
};

export type SearchResponse = {
  projects: ProjectResponse[];
  pager: PagerResponse;
};
