import type {
  PagerRequest,
  PagerResponse,
} from "../../../shared/types/dto/pager";
import type { Order } from "../../../shared/types/dto/order";

import type { TaskPriorityResponse } from "../../task-priorities/types/dto";
import type { TaskStatusResponse } from "../../task-statuses/types/dto";
import type { UserBaseResponse } from "../../users/types/dto";
import type { TimestampRange } from "../../../shared/composables/timestamps";

export type AddRequest = {
  summary: string;
  description: string | null;
  priority: {
    id: string;
  };
  status: {
    id: string;
  };
  tags: string[];
};

export type UpdateRequest = {
  id: string;
  summary: string;
  description: string | null;
  priority: {
    id: string;
  };
  status: {
    id: string;
  };
  startedAt: number | null;
  finishedAt: number | null;
  dueAt: number | null;
  tags: string[];
};

type SearchRequestFilter = {
  summary?: string;
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

export type TaskResponse = {
  id: string;
  projectId: string;
  slug: string;
  summary: string;
  description: string;
  priority: TaskPriorityResponse;
  status: TaskStatusResponse;
  createdAt: number;
  createdBy: UserBaseResponse;
  updatedAt?: number;
  startedAt?: number;
  finishedAt?: number;
  dueAt?: number;
  tags: string[];
  permissionsCount: number;
  attachmentsCount: number;
  notesCount: number;
  historyOperationsCount: number;
};

export type SearchResponse = {
  tasks: TaskResponse[];
  pager: PagerResponse;
};
