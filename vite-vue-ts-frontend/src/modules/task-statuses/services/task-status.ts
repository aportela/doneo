import { axiosInstance } from "../../../api/client";

import type {
  AddRequest,
  UpdateRequest,
  SearchRequest,
  TaskStatusResponse,
  SearchResponse,
} from "../types/dto";

export const taskStatusService = {
  async add(payload: AddRequest): Promise<TaskStatusResponse> {
    const { data } = await axiosInstance.post<TaskStatusResponse>(
      "/task-statuses",
      payload,
    );
    return data;
  },
  async update(payload: UpdateRequest): Promise<TaskStatusResponse> {
    const { data } = await axiosInstance.put<TaskStatusResponse>(
      "/task-statuses/" + payload.id,
      payload,
    );
    return data;
  },
  async delete(id: string): Promise<void> {
    await axiosInstance.delete<void>("/task-statuses/" + id);
  },
  async get(id: string): Promise<TaskStatusResponse> {
    const { data } = await axiosInstance.get<TaskStatusResponse>(
      "/task-statuses/" + id,
    );
    return data;
  },
  async searchBase(): Promise<SearchResponse> {
    const { data } = await axiosInstance.get<SearchResponse>(
      "/entities/task-statuses",
    );
    return data;
  },
  async search(payload: SearchRequest): Promise<SearchResponse> {
    const { data } = await axiosInstance.post<SearchResponse>(
      "/task-statuses/search",
      payload,
    );
    return data;
  },
};
