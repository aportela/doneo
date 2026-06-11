import { axiosInstance } from "../../../api/client";

import type { EmptyResponse, SearchResponse } from "../types/dto";

export const timerService = {
  async start(summary: string): Promise<EmptyResponse> {
    const payload = {
      summary: summary,
    };
    const { data } = await axiosInstance.post<EmptyResponse>(
      "/timers/",
      payload,
    );
    return data;
  },
  async stop(id: string): Promise<EmptyResponse> {
    const { data } = await axiosInstance.put<EmptyResponse>("/timers/" + id);
    return data;
  },
  async delete(id: string): Promise<EmptyResponse> {
    const { data } = await axiosInstance.delete<EmptyResponse>("/timers/" + id);
    return data;
  },
  async clear(): Promise<EmptyResponse> {
    const { data } = await axiosInstance.delete<EmptyResponse>("/timers/");
    return data;
  },
  async search(): Promise<SearchResponse> {
    const { data } = await axiosInstance.get<SearchResponse>("/timers/");
    return data;
  },
};
