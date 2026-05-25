import type {
  PagerRequest,
  PagerResponse,
} from "../../../shared/types/dto/pager";
import type { Order } from "../../../shared/types/dto/order";

import type { ProjectTypeResponse } from "../../project-types/types/dto";
import type { ProjectPriorityResponse } from "../../project-priorities/types/dto";
import type { ProjectStatusResponse } from "../../project-statuses/types/dto";
import type { UserBaseResponse } from "../../users/types/dto";

export type AddRequest = {
  key: string;
  summary: string;
};

export type UpdateRequest = {
  id: string;
  key: string;
  summary: string;
};

type SearchRequestFilter = {
  key?: string;
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
