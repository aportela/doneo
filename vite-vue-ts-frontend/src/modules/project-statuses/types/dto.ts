import type {
  PagerRequest,
  PagerResponse,
} from "../../../shared/types/dto/pager";
import type { Order } from "../../../shared/types/dto/order";

export type AddRequest = {
  name: string;
  HexColor: string;
};

export type UpdateRequest = {
  id: string;
  name: string;
  HexColor: string;
};

type SearchRequestFilter = {
  name?: string;
};

export type SearchRequest = {
  pager: PagerRequest;
  order: Order;
  filter?: SearchRequestFilter;
};

export type ProjectStatusResponse = {
  id: string;
  name: string;
  hexColor: string;
  index: number;
};

export type SearchResponse = {
  projectStatuses: ProjectStatusResponse[];
  pager: PagerResponse;
};
