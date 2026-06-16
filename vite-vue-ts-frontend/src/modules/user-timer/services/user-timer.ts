import { axiosInstance } from "../../../api/client";

import type { EmptyResponse, SearchResponse } from "../types/dto";

export const userTimerService = {
  async start(summary: string): Promise<EmptyResponse> {
    const payload = {
      summary: summary,
    };
    const { data } = await axiosInstance.post<EmptyResponse>(
      "/user-timers/",
      payload,
    );
    return data;
  },
  async stop(id: string): Promise<EmptyResponse> {
    const { data } = await axiosInstance.put<EmptyResponse>(
      "/user-timers/" + id,
    );
    return data;
  },
  async delete(id: string): Promise<EmptyResponse> {
    const { data } = await axiosInstance.delete<EmptyResponse>(
      "/user-timers/" + id,
    );
    return data;
  },
  async clear(): Promise<EmptyResponse> {
    const { data } = await axiosInstance.delete<EmptyResponse>("/user-timers/");
    return data;
  },
  async search(): Promise<SearchResponse> {
    const { data } = await axiosInstance.get<SearchResponse>("/user-timers/");
    return data;
  },
};
