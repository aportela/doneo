import { axiosInstance } from "../../../api/client";

import type { UpdateRequest, ProfileResponse } from "../types/dto";

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
};
