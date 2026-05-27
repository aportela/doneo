import type {
  PagerRequest,
  PagerResponse,
} from "../../../shared/types/dto/pager";
import type { Order } from "../../../shared/types/dto/order";

export type Permissions = {
  allowUpdateProject: boolean;
  allowDeleteProject: boolean;
  allowViewProject: boolean;
  allowAddTask: boolean;
  allowUpdateTask: boolean;
  allowDeleteTask: boolean;
  allowViewTask: boolean;
};

export const getDefaultPermissions = (): Permissions => {
  return {
    allowUpdateProject: false,
    allowDeleteProject: false,
    allowViewProject: false,
    allowAddTask: false,
    allowUpdateTask: false,
    allowDeleteTask: false,
    allowViewTask: false,
  };
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
};

export type SearchRequest = {
  pager: PagerRequest;
  order: Order;
  filter?: SearchRequestFilter;
};

export type RoleBaseResponse = {
  id: string;
  name: string;
};

// TODO: inherit
export type RoleResponse = {
  id: string;
  name: string;
  permissions: Permissions;
};

export type SearchBaseResponse = {
  roles: RoleBaseResponse[];
};

export type SearchResponse = {
  roles: RoleResponse[];
  pager: PagerResponse;
};
