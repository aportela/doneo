import type { UserBaseResponse } from "../../users/types/dto";

export type HistoryOperationResponse = {
  id: string;
  createdBy: UserBaseResponse;
  createdAt: number;
  operationType: number;
};

export type SearchResponse = {
  historyOperations: HistoryOperationResponse[];
};
