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
  name?: string;
  email?: string;
  isSuperUser?: boolean;
  page?: number;
  limit?: number;
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
};
