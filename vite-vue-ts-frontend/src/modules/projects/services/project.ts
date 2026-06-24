import { axiosInstance } from "../../../api/client";

import type {
  AddRequest,
  UpdateRequest,
  SearchRequest,
  ProjectResponse,
  SearchResponse,
  PatchRequest,
} from "../types/dto";

export const projectService = {
  async add(payload: AddRequest): Promise<ProjectResponse> {
    const { data } = await axiosInstance.post<ProjectResponse>(
      "/projects",
      payload,
    );
    return data;
  },
  async update(payload: UpdateRequest): Promise<ProjectResponse> {
    const { data } = await axiosInstance.put<ProjectResponse>(
      "/projects/" + payload.id,
      payload,
    );
    return data;
  },
  async patch(payload: PatchRequest): Promise<ProjectResponse> {
    const { data } = await axiosInstance.patch<ProjectResponse>(
      "/projects/" + payload.id,
      payload,
    );
    return data;
  },
  async delete(id: string): Promise<void> {
    await axiosInstance.delete<void>("/projects/" + id);
  },
  async get(id: string): Promise<ProjectResponse> {
    const { data } = await axiosInstance.get<ProjectResponse>(
      "/projects/" + id,
    );
    return data;
  },
  async search(payload: SearchRequest): Promise<SearchResponse> {
    const { data } = await axiosInstance.post<SearchResponse>(
      "/projects/search",
      payload,
    );
    return data;
  },
};
