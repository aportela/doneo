import { axiosInstance } from "../../../api/client";

import type {
  AddRequest,
  UpdateRequest,
  SearchRequest,
  UserResponse,
  SearchResponse,
} from "../types/dto";

export const userService = {
  async add(payload: AddRequest): Promise<UserResponse> {
    const { data } = await axiosInstance.post<UserResponse>("/users", payload);
    return data;
  },
  async update(payload: UpdateRequest): Promise<UserResponse> {
    const { data } = await axiosInstance.put<UserResponse>(
      "/users/" + payload.id,
      payload,
    );
    return data;
  },
  async delete(id: string): Promise<void> {
    await axiosInstance.delete<void>("/users/" + id);
  },
  async unDelete(id: string): Promise<void> {
    const params = {
      deletedAt: null,
    };
    await axiosInstance.patch<void>("/users/" + id, params);
  },
  async get(id: string): Promise<UserResponse> {
    const { data } = await axiosInstance.get<UserResponse>("/users/" + id);
    return data;
  },
  async search(payload: SearchRequest): Promise<SearchResponse> {
    const { data } = await axiosInstance.post<SearchResponse>(
      "/users/search",
      payload,
    );
    return data;
  },
};
