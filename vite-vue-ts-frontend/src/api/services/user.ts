import { axiosInstance } from "../client";
import type {
  AddResponseInterface,
  UpdateResponseInterface,
  GetResponseInterface,
  SearchResponseInterface,
} from "../types/dto/user";

import type { User } from "../models/user";

export const authService = {
  async add(user: User): Promise<AddResponseInterface> {
    const { data } = await axiosInstance.post<AddResponseInterface>(
      "/users",
      user,
    );
    return data;
  },
  async update(user: User): Promise<UpdateResponseInterface> {
    const { data } = await axiosInstance.post<UpdateResponseInterface>(
      "/users/" + user.id,
      user,
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
  async get(id: string): Promise<GetResponseInterface> {
    const { data } = await axiosInstance.get<GetResponseInterface>(
      "/users/" + id,
    );
    return data;
  },
  async search(): Promise<SearchResponseInterface> {
    const { data } = await axiosInstance.get<SearchResponseInterface>("/users");
    return data;
  },
};
