import type { UserBaseResponse } from "../../users/types/dto";

export type ProjectHistoryOperationResponse = {
  createdBy: UserBaseResponse;
  createdAt: number;
  operationType: number;
};

export type SearchResponse = {
  historyOperations: ProjectHistoryOperationResponse[];
};
