import type { PagerQuery, PagerResult } from "../../../shared/types/pager";
import type { Order } from "../../../shared/types/order";
// TODO: move to shared/types
import type { TimestampRange } from "../../../shared/composables/timestamps";

export interface UserPermissions {
  isSuperUser: boolean;
}

export interface AddRequest {
  name: string;
  email: string;
  password: string;
  permissions: UserPermissions;
}

export interface UpdateRequest {
  id: string;
  name: string;
  email: string;
  password?: string;
  permissions: UserPermissions;
}

export interface SearchRequest {
  pager: PagerQuery;
  order: Order;
  filter?: {
    type?: number;
    name?: string;
    email?: string;
    permissions?: {
      isSuperUser?: boolean;
    };
    createdAt?: TimestampRange;
    updatedAt?: TimestampRange;
    deletedAt?: TimestampRange;
  };
}

export interface UserBaseResponse {
  id: string;
  name: string;
}

export interface SearchBaseResponse {
  users: UserBaseResponse[];
}

export interface UserResponse {
  id: string;
  name: string;
  email: string;
  permissions: UserPermissions;
  createdAt: number;
  updatedAt: number | null;
  deletedAt: number | null;
}

export interface SearchResponse {
  users: UserResponse[];
  pager: PagerResult;
}
