import { axiosInstance } from "../../../api/client";

import type { SearchRequest, SearchResponse } from "../types/dto";

export const searchService = {
  async search(payload: SearchRequest): Promise<SearchResponse> {
    const { data } = await axiosInstance.post<SearchResponse>(
      "/search",
      payload,
    );
    return data;
  },
};
