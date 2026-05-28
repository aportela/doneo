import type { UserBaseResponse } from "../../users/types/dto";

export type AddRequest = {
  body: string;
};

export type UpdateRequest = {
  id: string;
  user: UserBaseResponse;
  createdAt: number;
  updatedAt: number | null;
  body: string;
};

export type NoteResponse = {
  id: string;
  user: UserBaseResponse;
  createdAt: number;
  updatedAt: number | null;
  body: string;
};

export type SearchResponse = {
  notes: NoteResponse[];
};
