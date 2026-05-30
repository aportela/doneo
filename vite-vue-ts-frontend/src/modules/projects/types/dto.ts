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
  key: string;
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
  key: string;
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

type SearchRequestFilter = {
  key?: string;
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

export type ProjectResponse = {
  id: string;
  key: string;
  summary: string;
  description: string;
  type: ProjectTypeResponse;
  priority: ProjectPriorityResponse;
  status: ProjectStatusResponse;
  createdAt: number;
  createdBy: UserBaseResponse;
  tasksCount: number;
  permissionsCount: number;
  attachmentsCount: number;
  notesCount: number;
  historyOperationsCount: number;
};

export type SearchResponse = {
  projects: ProjectResponse[];
  pager: PagerResponse;
};
