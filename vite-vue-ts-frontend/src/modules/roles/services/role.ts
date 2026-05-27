import { axiosInstance } from "../../../api/client";

import type {
  AddRequest,
  UpdateRequest,
  SearchRequest,
  RoleResponse,
  SearchBaseResponse,
  SearchResponse,
} from "../types/dto";

export const roleService = {
  async add(payload: AddRequest): Promise<RoleResponse> {
    const { data } = await axiosInstance.post<RoleResponse>("/roles", payload);
    return data;
  },
  async update(payload: UpdateRequest): Promise<RoleResponse> {
    const { data } = await axiosInstance.put<RoleResponse>(
      "/roles/" + payload.id,
      payload,
    );
    return data;
  },
  async delete(id: string): Promise<void> {
    await axiosInstance.delete<void>("/roles/" + id);
  },
  async get(id: string): Promise<RoleResponse> {
    const { data } = await axiosInstance.get<RoleResponse>("/roles/" + id);
    return data;
  },
  async searchBase(): Promise<SearchBaseResponse> {
    const { data } =
      await axiosInstance.get<SearchBaseResponse>("/entities/roles");
    return data;
  },
  async search(payload: SearchRequest): Promise<SearchResponse> {
    const { data } = await axiosInstance.post<SearchResponse>(
      "/roles/search",
      payload,
    );
    return data;
  },
};
