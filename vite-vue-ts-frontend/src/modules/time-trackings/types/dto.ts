import type { UserBaseResponse } from "../../users/types/dto";

export type AddRequest = {
  summary: string;
  spentTime: number;
};

export type UpdateRequest = {
  id: string;
  summary: string;
  spentTime: number;
};

export type TimeTrackingResponse = {
  id: string;
  createdBy: UserBaseResponse;
  createdAt: number;
  summary: string;
  spentTime: number;
};

export type SearchResponse = {
  timeTrackings: TimeTrackingResponse[];
};
