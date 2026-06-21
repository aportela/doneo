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
};
*/

interface UserPermissions {
  isSuperUser: boolean;
}

export type AddRequest = {
  name: string;
  email: string;
  permissions: UserPermissions;
};

export type UpdateRequest = {
  id: string;
  name: string;
  email: string;
  password?: string;
  permissions: UserPermissions;
};

interface SearchRequestUserPermissions {
  isSuperUser?: boolean;
}

type SearchRequestFilter = {
  type?: number;
  name?: string;
  email?: string;
  permissions?: SearchRequestUserPermissions;
  createdAt?: TimestampRange;
  updatedAt?: TimestampRange;
  deletedAt?: TimestampRange;
};

export type SearchRequest = {
  pager: PagerRequest;
  order: Order;
  filter?: SearchRequestFilter;
};

export type UserBaseResponse = {
  id: string;
  name: string;
};

export type UserResponse = {
  id: string;
  name: string;
  email: string;
  permissions: UserPermissions;
  createdAt: number;
  updatedAt: number | null;
  deletedAt: number | null;
};

export type SearchBaseResponse = {
  users: UserBaseResponse[];
};

export type SearchResponse = {
  users: UserResponse[];
  pager: PagerResponse;
};
