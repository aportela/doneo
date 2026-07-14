import { axiosInstance } from "../../../composables/axios";

import type { CurrentProjectsResponse } from "../types/dto";

export const navigationMenuService = {
  async getCurrentProjects(): Promise<CurrentProjectsResponse> {
    const { data } = await axiosInstance.get<CurrentProjectsResponse>(
      "/menu/current_projects",
    );
    return data;
  },
};
