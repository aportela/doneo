import type { PagerQuery, PagerResult } from "../../../shared/types/pager";
import type { Order } from "../../../shared/types/order";

export type AddRequest = {
  name: string;
  hexColor: string;
  index: number;
};

export type UpdateRequest = {
  id: string;
  name: string;
  hexColor: string;
  index: number;
};

type SearchRequestFilter = {
  name?: string;
};

export type SearchRequest = {
  pager: PagerQuery;
  order: Order;
  filter?: SearchRequestFilter;
};

export type TaskPriorityResponse = {
  id: string;
  name: string;
  hexColor: string;
  index: number;
};

export type SearchResponse = {
  taskPriorities: TaskPriorityResponse[];
  pager: PagerResult;
};
