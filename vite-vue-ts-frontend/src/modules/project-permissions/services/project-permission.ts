import { axiosInstance } from "../../../api/client";

import type {
  AddRequest,
  ProjectPermissionResponse,
  SearchResponse,
} from "../types/dto";

export const projectPermissionService = {
  async add(
    projectId: string,
    payload: AddRequest,
  ): Promise<ProjectPermissionResponse> {
    const { data } = await axiosInstance.post<ProjectPermissionResponse>(
      "/projects/" + projectId + "/permissions",
      payload,
    );
    return data;
  },
  async delete(projectId: string, permissionId: string): Promise<void> {
    await axiosInstance.delete<void>(
      "/projects/" + projectId + "/permissions/" + permissionId,
    );
  },
  async search(projectId: string): Promise<SearchResponse> {
    const { data } = await axiosInstance.get<SearchResponse>(
      "/projects/" + projectId + "/permissions",
    );
    return data;
  },
};
