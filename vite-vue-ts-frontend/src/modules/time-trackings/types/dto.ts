import type { UserBaseResponse } from "../../users/types/dto";

export type AddRequest = {
  summary: string;
  totalSeconds: number;
};

export type UpdateRequest = {
  id: string;
  summary: string;
  totalSeconds: number;
};

export type TimeTrackingResponse = {
  id: string;
  createdBy: UserBaseResponse;
  createdAt: number;
  summary: string;
  totalSeconds: number;
};

export type SearchResponse = {
  timeTrackings: TimeTrackingResponse[];
};
