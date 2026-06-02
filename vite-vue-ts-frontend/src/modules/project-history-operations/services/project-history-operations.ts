import { axiosInstance } from "../../../api/client";

import type { SearchResponse } from "../types/dto";

export const projectHistoryOperationsService = {
  async getProjectHistoryOperations(
    projectId: string,
  ): Promise<SearchResponse> {
    const { data } = await axiosInstance.get<SearchResponse>(
      "/projects/" + projectId + "/history_operations",
    );
    return data;
  },
};
