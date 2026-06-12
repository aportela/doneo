import { axiosInstance } from "../../../api/client";

import type { SearchResponse } from "../types/dto";

export const historyOperationsService = {
  async getProjectHistoryOperations(
    projectId: string,
  ): Promise<SearchResponse> {
    const { data } = await axiosInstance.get<SearchResponse>(
      "/projects/" + projectId + "/history_operations",
    );
    return data;
  },
  async getTaskHistoryOperations(
    projectId: string,
    taskId: string,
  ): Promise<SearchResponse> {
    const { data } = await axiosInstance.get<SearchResponse>(
      "/projects/" + projectId + "/tasks/" + taskId + "/history_operations",
    );
    return data;
  },
};
