import { axiosInstance } from "../../../api/client";

import type {
  AddRequest,
  UpdateRequest,
  PageResponse,
  SearchResponse,
} from "../types/dto";

export const pageService = {
  async addProjectPage(
    projectId: string,
    payload: AddRequest,
  ): Promise<PageResponse> {
    const { data } = await axiosInstance.post<PageResponse>(
      "/projects/" + projectId + "/pages",
      payload,
    );
    return data;
  },
  async updateProjectPage(
    projectId: string,
    pageId: string,
    payload: UpdateRequest,
  ): Promise<PageResponse> {
    const { data } = await axiosInstance.put<PageResponse>(
      "/projects/" + projectId + "/pages/" + pageId,
      payload,
    );
    return data;
  },
  async deleteProjectPage(projectId: string, pageId: string): Promise<void> {
    await axiosInstance.delete<void>(
      "/projects/" + projectId + "/pages/" + pageId,
    );
  },
  async getProjectPage(
    projectId: string,
    pageId: string,
  ): Promise<PageResponse> {
    const { data } = await axiosInstance.get<PageResponse>(
      "/projects/" + projectId + "/pages/" + pageId,
    );
    return data;
  },
  async getProjectPages(projectId: string): Promise<SearchResponse> {
    const { data } = await axiosInstance.get<SearchResponse>(
      "/projects/" + projectId + "/pages",
    );
    return data;
  },
};
