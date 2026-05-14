import type {
  PagerRequest,
  PagerResponse,
} from "../../../shared/types/dto/pager";
import type { Order } from "../../../shared/types/dto/order";
import type { TimestampRange } from "../../../shared/composables/timestamps";

/*
export type UserRequest = {
  id: string;
  name: string;
  email: string;
  isSuperUser: boolean;
  createdAt: number;
  updatedAt: number | null;
  deletedAt: number | null;
  avatarURL: string;
};
*/

export type AddRequest = {
  name: string;
  email: string;
  isSuperUser: boolean;
};

export type UpdateRequest = {
  id: string;
  name: string;
  email: string;
  password?: string;
  isSuperUser: boolean;
};

type SearchRequestFilter = {
  type?: number;
  name?: string;
  email?: string;
  createdAt?: TimestampRange;
  updatedAt?: TimestampRange;
  deletedAt?: TimestampRange;
};

export type SearchRequest = {
  pager: PagerRequest;
  order: Order;
  filter?: SearchRequestFilter;
};

export type UserResponse = {
  id: string;
  name: string;
  email: string;
  isSuperUser: boolean;
  createdAt: number;
  updatedAt: number | null;
  deletedAt: number | null;
  avatarUrl: string | null;
};

export type SearchResponse = {
  users: UserResponse[];
  pager: PagerResponse;
};
