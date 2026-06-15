import { axiosInstance } from "../../../api/client";

import type {
  AddRequest,
  UpdateRequest,
  SearchRequest,
  TaskResponse,
  SearchResponse,
} from "../types/dto";

export const taskService = {
  async add(projectId: string, payload: AddRequest): Promise<TaskResponse> {
    const { data } = await axiosInstance.post<TaskResponse>(
      "/projects/" + projectId + "/tasks",
      payload,
    );
    return data;
  },
  async update(
    projectId: string,
    payload: UpdateRequest,
  ): Promise<TaskResponse> {
    const { data } = await axiosInstance.put<TaskResponse>(
      "/projects/" + projectId + "/tasks/" + payload.id,
      payload,
    );
    return data;
  },
  async delete(projectId: string, id: string): Promise<void> {
    await axiosInstance.delete<void>("/projects/" + projectId + "/tasks/" + id);
  },
  async get(projectId: string, id: string): Promise<TaskResponse> {
    const { data } = await axiosInstance.get<TaskResponse>(
      "/projects/" + projectId + "/tasks/" + id,
    );
    return data;
  },
  async search(
    projectId: string | null,
    payload: SearchRequest,
  ): Promise<SearchResponse> {
    const { data } = await axiosInstance.post<SearchResponse>(
      projectId ? "/projects/" + projectId + "/tasks/search" : "/tasks/search",
      payload,
    );
    return data;
  },
};
