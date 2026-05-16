import type {
  PagerRequest,
  PagerResponse,
} from "../../../shared/types/dto/pager";
import type { Order } from "../../../shared/types/dto/order";

export type Permissions = {
  allowCreate: boolean;
  allowUpdate: boolean;
  allowDelete: boolean;
  allowView: boolean;
  allowList: boolean;
  allowExecute: boolean;
};

export type AddRequest = {
  name: string;
  permissions: Permissions;
};

export type UpdateRequest = {
  id: string;
  name: string;
  permissions: Permissions;
};

type SearchRequestFilter = {
  name?: string;
  permissions?: {
    allowCreate?: boolean;
    allowUpdate?: boolean;
    allowDelete?: boolean;
    allowView?: boolean;
    allowList?: boolean;
    allowExecute?: boolean;
  };
};

export type SearchRequest = {
  pager: PagerRequest;
  order: Order;
  filter?: SearchRequestFilter;
};

export type RoleResponse = {
  id: string;
  name: string;
  permissions: Permissions;
};

export type SearchResponse = {
  roles: RoleResponse[];
  pager: PagerResponse;
};
