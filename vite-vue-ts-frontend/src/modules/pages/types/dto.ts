import type { UserBaseResponse } from "../../users/types/dto";

export type AddRequest = {
  title: string;
};

export type UpdateRequest = {
  id: string;
  title: string;
  body: string;
};

export type PageResponse = {
  id: string;
  createdBy: UserBaseResponse;
  createdAt: number;
  updatedAt: number | null;
  title: string;
  body: string;
};

export type SearchResponse = {
  pages: PageResponse[];
};
