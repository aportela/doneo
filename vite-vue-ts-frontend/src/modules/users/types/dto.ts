import type {
  PagerRequest,
  PagerResponse,
} from "../../../shared/types/dto/pager";
import type { Order } from "../../../shared/types/dto/order";

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
  user: {
    name: string;
    email: string;
    isSuperUser: boolean;
  };
};

export type UpdateRequest = {
  user: {
    id: string;
    name: string;
    email: string;
    isSuperUser: boolean;
  };
};

export type SearchRequest = {
  pager: PagerRequest;
  order: Order;
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

export type AddResponse = {
  user: UserResponse;
};

export type UpdateResponse = {
  user: UserResponse;
};

export type GetResponse = {
  user: UserResponse;
};

export type SearchResponse = {
  users: UserResponse[];
  pager: PagerResponse;
};
