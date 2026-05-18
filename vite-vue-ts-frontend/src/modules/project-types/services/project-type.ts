import { axiosInstance } from "../../../api/client";

import type {
  AddRequest,
  UpdateRequest,
  SearchRequest,
  ProjectTypeResponse,
  SearchResponse,
} from "../types/dto";

export const projectTypeService = {
  async add(payload: AddRequest): Promise<ProjectTypeResponse> {
    const { data } = await axiosInstance.post<ProjectTypeResponse>(
      "/project-types",
      payload,
    );
    return data;
  },
  async update(payload: UpdateRequest): Promise<ProjectTypeResponse> {
    const { data } = await axiosInstance.put<ProjectTypeResponse>(
      "/project-types/" + payload.id,
      payload,
    );
    return data;
  },
  async delete(id: string): Promise<void> {
    await axiosInstance.delete<void>("/project-types/" + id);
  },
  async get(id: string): Promise<ProjectTypeResponse> {
    const { data } = await axiosInstance.get<ProjectTypeResponse>(
      "/project-types/" + id,
    );
    return data;
  },
  async search(payload: SearchRequest): Promise<SearchResponse> {
    const { data } = await axiosInstance.post<SearchResponse>(
      "/project-types/search",
      payload,
    );
    return data;
  },
};
