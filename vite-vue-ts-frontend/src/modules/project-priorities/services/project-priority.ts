import { axiosInstance } from "../../../api/client";

import type {
  AddRequest,
  UpdateRequest,
  SearchRequest,
  ProjectPriorityResponse,
  SearchResponse,
} from "../types/dto";

export const projectPriorityService = {
  async add(payload: AddRequest): Promise<ProjectPriorityResponse> {
    const { data } = await axiosInstance.post<ProjectPriorityResponse>(
      "/project-priorities",
      payload,
    );
    return data;
  },
  async update(payload: UpdateRequest): Promise<ProjectPriorityResponse> {
    const { data } = await axiosInstance.put<ProjectPriorityResponse>(
      "/project-priorities/" + payload.id,
      payload,
    );
    return data;
  },
  async delete(id: string): Promise<void> {
    await axiosInstance.delete<void>("/project-priorities/" + id);
  },
  async get(id: string): Promise<ProjectPriorityResponse> {
    const { data } = await axiosInstance.get<ProjectPriorityResponse>(
      "/project-priorities/" + id,
    );
    return data;
  },
  async searchBase(): Promise<SearchResponse> {
    const { data } = await axiosInstance.get<SearchResponse>(
      "/entities/project-priorities",
    );
    return data;
  },
  async search(payload: SearchRequest): Promise<SearchResponse> {
    const { data } = await axiosInstance.post<SearchResponse>(
      "/project-priorities/search",
      payload,
    );
    return data;
  },
};
