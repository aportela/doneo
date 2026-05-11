import { axiosInstance } from "../client";
import type {
  AddRequest,
  AddResponse,
  UpdateRequest,
  SearchRequest,
  UpdateResponse,
  GetResponse,
  SearchResponse,
} from "../types/dto/user";

export const userService = {
  async add(payload: AddRequest): Promise<AddResponse> {
    const { data } = await axiosInstance.post<AddResponse>("/users", payload);
    return data;
  },
  async update(payload: UpdateRequest): Promise<UpdateResponse> {
    const { data } = await axiosInstance.post<UpdateResponse>(
      "/users/" + payload.user.id,
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
  async get(id: string): Promise<GetResponse> {
    const { data } = await axiosInstance.get<GetResponse>("/users/" + id);
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
