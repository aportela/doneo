import type {
  PagerRequest,
  PagerResponse,
} from "../../../shared/types/dto/pager";
import type { Order } from "../../../shared/types/dto/order";

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
  pager: PagerRequest;
  order: Order;
  filter?: SearchRequestFilter;
};

export type TaskStatusResponse = {
  id: string;
  name: string;
  hexColor: string;
  index: number;
};

export type SearchResponse = {
  taskStatuses: TaskStatusResponse[];
  pager: PagerResponse;
};
