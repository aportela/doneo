import { axiosInstance } from "../../../api/client";

import type {
  AddRequest,
  UpdateRequest,
  SearchRequest,
  ProjectStatusResponse,
  SearchResponse,
} from "../types/dto";

export const projectStatusService = {
  async add(payload: AddRequest): Promise<ProjectStatusResponse> {
    const { data } = await axiosInstance.post<ProjectStatusResponse>(
      "/project-statuses",
      payload,
    );
    return data;
  },
  async update(payload: UpdateRequest): Promise<ProjectStatusResponse> {
    const { data } = await axiosInstance.put<ProjectStatusResponse>(
      "/project-statuses/" + payload.id,
      payload,
    );
    return data;
  },
  async delete(id: string): Promise<void> {
    await axiosInstance.delete<void>("/project-statuses/" + id);
  },
  async get(id: string): Promise<ProjectStatusResponse> {
    const { data } = await axiosInstance.get<ProjectStatusResponse>(
      "/project-statuses/" + id,
    );
    return data;
  },
  async searchBase(): Promise<SearchResponse> {
    const { data } = await axiosInstance.get<SearchResponse>(
      "/entities/project-statuses",
    );
    return data;
  },
  async search(payload: SearchRequest): Promise<SearchResponse> {
    const { data } = await axiosInstance.post<SearchResponse>(
      "/project-statuses/search",
      payload,
    );
    return data;
  },
};
