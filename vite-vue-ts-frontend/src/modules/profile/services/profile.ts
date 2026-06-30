import { axiosInstance } from "../../../api/client";

import type {
  UpdateRequest,
  SaveAvatarRequest,
  ProfileResponse,
  EmptyResponse,
} from "../types/dto";

export const profileService = {
  async update(payload: UpdateRequest): Promise<ProfileResponse> {
    const { data } = await axiosInstance.put<ProfileResponse>(
      "/profile/",
      payload,
    );
    return data;
  },
  async get(): Promise<ProfileResponse> {
    const { data } = await axiosInstance.get<ProfileResponse>("/profile/");
    return data;
  },
  async saveAvatar(payload: SaveAvatarRequest): Promise<EmptyResponse> {
    const { data } = await axiosInstance.put<ProfileResponse>(
      "/profile/avatar/",
      payload,
    );
    return data;
  },
  async deleteAvatar(): Promise<EmptyResponse> {
    const { data } =
      await axiosInstance.delete<ProfileResponse>("/profile/avatar/");
    return data;
  },
};
