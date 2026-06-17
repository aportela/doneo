import { axiosInstance } from "../../../api/client";

import type {
  AddRequest,
  UpdateRequest,
  SearchRequest,
  TaskPriorityResponse,
  SearchResponse,
} from "../types/dto";

export const taskPriorityService = {
  async add(payload: AddRequest): Promise<TaskPriorityResponse> {
    const { data } = await axiosInstance.post<TaskPriorityResponse>(
      "/task-priorities",
      payload,
    );
    return data;
  },
  async update(payload: UpdateRequest): Promise<TaskPriorityResponse> {
    const { data } = await axiosInstance.put<TaskPriorityResponse>(
      "/task-priorities/" + payload.id,
      payload,
    );
    return data;
  },
  async delete(id: string): Promise<void> {
    await axiosInstance.delete<void>("/task-priorities/" + id);
  },
  async get(id: string): Promise<TaskPriorityResponse> {
    const { data } = await axiosInstance.get<TaskPriorityResponse>(
      "/task-priorities/" + id,
    );
    return data;
  },
  async searchBase(): Promise<SearchResponse> {
    const { data } = await axiosInstance.get<SearchResponse>(
      "/entities/task-priorities",
    );
    return data;
  },
  async search(payload: SearchRequest): Promise<SearchResponse> {
    const { data } = await axiosInstance.post<SearchResponse>(
      "/task-priorities/search",
      payload,
    );
    return data;
  },
};
